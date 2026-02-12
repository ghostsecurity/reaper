package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/ghostsecurity/reaper/internal/storage"
)

type IPCServer struct {
	listener net.Listener
	store    storage.Store
	shutdown chan struct{}
}

func NewIPCServer(dataDir string, store storage.Store, shutdown chan struct{}) (*IPCServer, error) {
	sockPath := filepath.Join(dataDir, "reaper.sock")

	// Remove stale socket
	os.Remove(sockPath)

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		return nil, fmt.Errorf("listening on socket: %w", err)
	}

	return &IPCServer{
		listener: listener,
		store:    store,
		shutdown: shutdown,
	}, nil
}

func (s *IPCServer) Serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return
			default:
				continue
			}
		}
		go s.handleConn(conn)
	}
}

func (s *IPCServer) Close() error {
	return s.listener.Close()
}

func (s *IPCServer) handleConn(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 10*1024*1024), 10*1024*1024)
	if !scanner.Scan() {
		return
	}

	var req Request
	if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
		writeResponse(conn, Response{Error: "invalid request"})
		return
	}

	resp := s.route(req)
	writeResponse(conn, resp)
}

func (s *IPCServer) route(req Request) Response {
	switch req.Command {
	case "logs":
		return s.handleLogs(req.Params)
	case "search":
		return s.handleSearch(req.Params)
	case "get", "req", "res":
		return s.handleGet(req.Command, req.Params)
	case "tail":
		return s.handleTail(req.Params)
	case "clear":
		return s.handleClear()
	case "shutdown":
		return s.handleShutdown()
	case "ping":
		return Response{OK: true}
	default:
		return Response{Error: fmt.Sprintf("unknown command: %s", req.Command)}
	}
}

func (s *IPCServer) handleLogs(params json.RawMessage) Response {
	var p LogsParams
	if len(params) > 0 {
		if err := json.Unmarshal(params, &p); err != nil {
			return Response{Error: "invalid params"}
		}
	}
	if p.Limit <= 0 {
		p.Limit = 50
	}

	entries, err := s.store.List(p.Limit, p.Offset)
	if err != nil {
		return Response{Error: err.Error()}
	}

	data, _ := json.Marshal(entries)
	return Response{OK: true, Data: data}
}

func (s *IPCServer) handleSearch(params json.RawMessage) Response {
	var p SearchRequestParams
	if len(params) > 0 {
		if err := json.Unmarshal(params, &p); err != nil {
			return Response{Error: "invalid params"}
		}
	}

	entries, err := s.store.Search(storage.SearchParams{
		Method:  p.Method,
		Host:    p.Host,
		Domains: p.Domains,
		Path:    p.Path,
		Status:  p.Status,
		Limit:   p.Limit,
		Offset:  p.Offset,
	})
	if err != nil {
		return Response{Error: err.Error()}
	}

	data, _ := json.Marshal(entries)
	return Response{OK: true, Data: data}
}

func (s *IPCServer) handleGet(command string, params json.RawMessage) Response {
	var p GetParams
	if err := json.Unmarshal(params, &p); err != nil {
		return Response{Error: "invalid params"}
	}

	entry, err := s.store.Get(p.ID)
	if err != nil {
		return Response{Error: err.Error()}
	}

	type getResponse struct {
		Command string         `json:"command"`
		Entry   *storage.Entry `json:"entry"`
	}

	data, _ := json.Marshal(getResponse{Command: command, Entry: entry})
	return Response{OK: true, Data: data}
}

func (s *IPCServer) handleTail(params json.RawMessage) Response {
	var p TailParams
	if len(params) > 0 {
		if err := json.Unmarshal(params, &p); err != nil {
			return Response{Error: "invalid params"}
		}
	}

	entries, err := s.store.ListAfter(p.AfterID, p.Limit)
	if err != nil {
		return Response{Error: err.Error()}
	}

	data, _ := json.Marshal(entries)
	return Response{OK: true, Data: data}
}

func (s *IPCServer) handleClear() Response {
	if err := s.store.Clear(); err != nil {
		return Response{Error: err.Error()}
	}
	return Response{OK: true}
}

func (s *IPCServer) handleShutdown() Response {
	go func() {
		close(s.shutdown)
	}()
	return Response{OK: true}
}

func writeResponse(conn net.Conn, resp Response) {
	data, _ := json.Marshal(resp)
	data = append(data, '\n')
	_, _ = conn.Write(data)
}
