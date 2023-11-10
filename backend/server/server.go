package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/net/context"

	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/server/api"
	"github.com/ghostsecurity/reaper/backend/settings"
)

type Server struct {
	staticFS         fs.FS
	upgrade          websocket.Upgrader
	connMu           sync.Mutex
	conns            []*threadSafeWriter
	api              *api.API
	logger           *log.Logger
	ctx              context.Context
	settingsProvider *settings.Provider
}

func New(ctx context.Context, staticFS fs.FS, logger *log.Logger, settingsProvider *settings.Provider) *Server {
	srv := &Server{
		staticFS: staticFS,
		upgrade: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		logger:           logger,
		ctx:              ctx,
		settingsProvider: settingsProvider,
	}
	srv.api = api.New(ctx, logger, settingsProvider, srv.TriggerEvent)
	return srv

}

type dir string

const (
	dirIn  dir = "<-"
	dirOut dir = "->"
)

func (s *Server) log(d dir, format string, args ...interface{}) {
	s.logger.Debugf("%s %s", d, fmt.Sprintf(format, args...))
}

func (s *Server) Start(addr string) error {

	dist, err := fs.Sub(s.staticFS, "dist")
	if err != nil {
		return fmt.Errorf("failed to create dist file system: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(dist)))
	mux.HandleFunc("/ws/", s.websocketHandler)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.logger.Infof("Listening on %s...", listener.Addr())

	srv := &http.Server{
		BaseContext: func(ln net.Listener) context.Context {
			return s.ctx
		},
		Handler: mux,
	}

	var closing bool
	go func() {
		<-s.ctx.Done()
		closing = true
		s.logger.Infof("Signal received, shutting down...")
		_ = listener.Close()
		_ = srv.Close()
	}()

	s.logger.Infof("Starting server...")
	err = srv.Serve(listener)

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		if closing {
			s.logger.Infof("Server closed")
			return s.ctx.Err()
		}
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

type MessageType uint8

const (
	MessageTypeUnknown MessageType = iota
	MessageTypePing
	MessageTypePong
	MessageTypeSubscribe
	MessageTypeNotify
	MessageTypeMethod
	MessageTypeResult
	MessageTypeFailure
	MessageTypeError
)

func (mt MessageType) String() string {
	switch mt {
	case MessageTypeUnknown:
		return "unknown"
	case MessageTypePing:
		return "ping"
	case MessageTypePong:
		return "pong"
	case MessageTypeSubscribe:
		return "subscribe"
	case MessageTypeNotify:
		return "notify"
	case MessageTypeMethod:
		return "method"
	case MessageTypeResult:
		return "result"
	case MessageTypeFailure:
		return "failure"
	case MessageTypeError:
		return "error"
	default:
		return fmt.Sprintf("unknown(%d)", mt)
	}
}

type websocketMessage struct {
	MessageType MessageType `json:"messageType"`
	Identifier  string      `json:"identifier"`
	Args        []string    `json:"args"` // json for each arg
	Sender      string      `json:"sender"`
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	unsafeConn, err := s.upgrade.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("failed to upgrade websocket connection: %s", err)
		return
	}
	// Close the connection when the for-loop operation is finished.
	defer func() { _ = unsafeConn.Close() }()

	conn := newThreadSafeWriter(unsafeConn, s.log)
	s.connMu.Lock()
	s.conns = append(s.conns, conn)
	s.connMu.Unlock()
	defer func() {
		s.connMu.Lock()
		defer s.connMu.Unlock()
		for i, c := range s.conns {
			if c == conn {
				s.conns = append(s.conns[:i], s.conns[i+1:]...)
				return
			}
		}
	}()

	s.log(dirOut, "ping:server")
	if err := conn.WriteJSON(&websocketMessage{
		MessageType: MessageTypePing,
		Identifier:  "server",
	}); err != nil {
		s.logger.Errorf("failed to write ping: %s", err)
		return
	}

	var message websocketMessage
	for {
		// the first message is "connected"
		_, raw, err := conn.ReadMessage()
		if err != nil {
			s.logger.Errorf("failed to read message: %s", err)
			break
		} else if err := json.Unmarshal(raw, &message); err != nil {
			s.logger.Errorf("failed to unmarshal message: %s", err)
			break
		}
		s.log(dirIn, "%s:%s", message.MessageType, message.Identifier)
		reply := s.handleMessage(message, conn)
		if reply != nil {
			s.log(dirOut, "%s:%s", reply.MessageType, reply.Identifier)
			if err := conn.WriteJSON(reply); err != nil {
				s.logger.Errorf("failed to write message: %s", err)
				return
			}
		}
	}
}

func (s *Server) handleMessage(message websocketMessage, conn *threadSafeWriter) *websocketMessage {
	switch message.MessageType {
	case MessageTypeUnknown:
		return createErrorMessage("Unspecified message type: %s", message.MessageType)
	case MessageTypePing:
		return &websocketMessage{
			MessageType: MessageTypePong,
			Identifier:  message.Identifier,
		}
	case MessageTypePong:
		return nil
	case MessageTypeSubscribe:
		s.subscribeEvent(message.Identifier)
		return nil
	case MessageTypeNotify:
		return createErrorMessage("Unexpected notify: %s", message.Identifier)
	case MessageTypeMethod:
		output, err := s.callMethod(message.Identifier, message.Args)
		if err != nil {
			errDat, _ := json.Marshal(err.Error())
			return &websocketMessage{
				MessageType: MessageTypeFailure,
				Identifier:  message.Identifier,
				Args:        []string{string(errDat)},
				Sender:      message.Sender,
			}
		}
		return &websocketMessage{
			MessageType: MessageTypeResult,
			Identifier:  message.Identifier,
			Args:        output,
			Sender:      message.Sender,
		}
	case MessageTypeResult:
		return createErrorMessage("Unexpected result: %s", message.Identifier)
	case MessageTypeError:
		// TODO: handle client error?
		return nil
	default:
		return createErrorMessage("Unknown message type: %s", message.MessageType)
	}
}

func createErrorMessage(msg string, args ...interface{}) *websocketMessage {
	dat, _ := json.Marshal(fmt.Sprintf(msg, args...))
	return &websocketMessage{
		MessageType: MessageTypeError,
		Identifier:  string(dat),
	}
}

func (s *Server) subscribeEvent(event string) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	for _, conn := range s.conns {
		conn.subscribeEvent(event)
	}
}

func (s *Server) TriggerEvent(event string, args ...interface{}) error {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	for _, conn := range s.conns {
		if err := conn.NotifyEvent(event, args...); err != nil {
			return fmt.Errorf("failed to trigger event: %w", err)
		}
	}
	return nil
}

func (s *Server) callMethod(method string, args []string) ([]string, error) {
	f := reflect.ValueOf(s.api).MethodByName(method)
	if f.IsZero() {
		return nil, fmt.Errorf("method not found")
	}
	inputs := make([]reflect.Value, f.Type().NumIn())
	for i, arg := range args {
		v := reflect.New(f.Type().In(i))
		// reflected pointer
		if err := json.Unmarshal([]byte(arg), v.Interface()); err != nil {
			return nil, fmt.Errorf("failed to unmarshal input: %w", err)
		}
		inputs[i] = v.Elem()
	}
	outputs := f.Call(inputs)
	var rawOut []string
	for i, output := range outputs {
		if f.Type().Out(i).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			if !output.IsNil() {
				return nil, fmt.Errorf("method returned error: %w", output.Interface().(error))
			}
			continue
		}
		data, err := json.Marshal(output.Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal output: %w", err)
		}
		rawOut = append(rawOut, string(data))
	}
	return rawOut, nil
}
