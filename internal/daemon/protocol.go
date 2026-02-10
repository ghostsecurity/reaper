package daemon

import "encoding/json"

type Request struct {
	Command string          `json:"command"` // "logs", "search", "get", "req", "res", "shutdown"
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

type LogsParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SearchRequestParams struct {
	Method  string   `json:"method,omitempty"`
	Host    string   `json:"host,omitempty"`
	Domains []string `json:"domains,omitempty"`
	Path    string   `json:"path,omitempty"`
	Status  int      `json:"status,omitempty"`
	Limit   int      `json:"limit,omitempty"`
	Offset  int      `json:"offset,omitempty"`
}

type GetParams struct {
	ID int64 `json:"id"`
}
