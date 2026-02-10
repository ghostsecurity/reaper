package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"time"
)

type Client struct {
	sockPath string
}

func NewClient(dataDir string) *Client {
	return &Client{
		sockPath: filepath.Join(dataDir, "reaper.sock"),
	}
}

func (c *Client) Send(req Request) (*Response, error) {
	conn, err := net.DialTimeout("unix", c.sockPath, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connecting to daemon: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}
	data = append(data, '\n')

	if _, err := conn.Write(data); err != nil {
		return nil, fmt.Errorf("writing request: %w", err)
	}

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 10*1024*1024), 10*1024*1024)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}
		return nil, fmt.Errorf("no response from daemon")
	}

	var resp Response
	if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp, nil
}

func WaitForSocket(dataDir string) error {
	sockPath := filepath.Join(dataDir, "reaper.sock")
	for i := 0; i < 30; i++ {
		conn, err := net.DialTimeout("unix", sockPath, time.Second)
		if err == nil {
			// Send ping
			data, _ := json.Marshal(Request{Command: "ping"})
			data = append(data, '\n')
			conn.Write(data)
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for daemon socket")
}
