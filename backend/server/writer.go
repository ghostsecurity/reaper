package server

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
	events map[string]struct{}
	log    func(d dir, format string, args ...interface{})
}

func newThreadSafeWriter(conn *websocket.Conn, log func(d dir, format string, args ...interface{})) *threadSafeWriter {
	return &threadSafeWriter{
		Conn:   conn,
		events: make(map[string]struct{}),
		log:    log,
	}
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()
	return t.Conn.WriteJSON(v)
}

func (t *threadSafeWriter) subscribeEvent(event string) {
	t.Lock()
	defer t.Unlock()
	t.events[event] = struct{}{}
}

func (t *threadSafeWriter) NotifyEvent(event string, args ...interface{}) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.events[event]; !ok {
		return nil
	}
	t.log(dirOut, "event:%s", event)
	values := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == nil {
			continue
		}
		value, err := json.Marshal(arg)
		if err != nil {
			return fmt.Errorf("failed to marshal event data: %w", err)
		}
		values = append(values, string(value))
	}
	return t.Conn.WriteJSON(websocketMessage{
		MessageType: MessageTypeNotify,
		Identifier:  event,
		Args:        values,
	})
}
