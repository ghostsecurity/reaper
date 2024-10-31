package replay

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ghostsecurity/reaper/internal/types"
)

func Do(ctx context.Context, req *types.ReplayInput) (*types.ReplayResponse, error) {
	var request *http.Request
	var body io.Reader

	// setup HTTP client with a sensible timeout
	t := http.DefaultTransport.(*http.Transport).Clone()
	client := &http.Client{
		// TODO: make timeout configurable
		Timeout:   2 * time.Second,
		Transport: t,
	}

	if req.Body != "" {
		body = bytes.NewBufferString(req.Body)
	}

	// new request
	request, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return nil, err
	}

	// set headers
	for key, value := range parseHeaderText(req.Headers) {
		request.Header.Set(key, value)
	}

	// send the request
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	responseBody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	replayResponse := types.ReplayResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       responseBody,
	}

	return &replayResponse, nil
}

// parseHeaderText parses a newline delimited header string into a map
func parseHeaderText(header string) map[string]string {
	m := make(map[string]string)
	headers := strings.Split(header, "\n")

	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			key := strings.Trim(parts[0], " ")   // header key
			value := strings.Trim(parts[1], " ") // header value
			if !shouldIgnoreHeader(key, value) {
				m[key] = value
			}
		}
	}

	return m
}

func shouldIgnoreHeader(key, val string) bool {
	// ignore connection: keep-alive
	// we can't re-use a connection for a replay
	if strings.ToLower(key) == "connection" && strings.ToLower(val) == "keep-alive" {
		return true
	}

	return false
}
