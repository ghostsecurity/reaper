package storage

import (
	"net/http"
	"time"
)

type Entry struct {
	ID              int64
	Method          string
	Scheme          string // "http" or "https"
	Host            string
	Path            string
	Query           string
	RequestHeaders  http.Header
	RequestBody     []byte
	StatusCode      int
	ResponseHeaders http.Header
	ResponseBody    []byte
	Timestamp       time.Time
	DurationMs      int64
}

type SearchParams struct {
	Method  string   // exact match
	Host    string   // supports glob wildcard (*.domain.com)
	Domains []string // suffix match
	Path    string   // prefix or glob
	Status  int      // exact match
	Limit   int
	Offset  int
}
