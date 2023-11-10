// Package packaging provides functions to convert to and from http.Request/Response and json marshallable equivalents.
package packaging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HttpRequest struct {
	Method      string        `json:"method"`
	URL         string        `json:"url"`
	Host        string        `json:"host"`
	Path        string        `json:"path"`
	QueryString string        `json:"query_string"`
	Scheme      string        `json:"scheme"`
	Body        string        `json:"body"`
	ID          string        `json:"id"`
	LocalID     int64         `json:"local_id"`
	Headers     []KeyValue    `json:"headers"`
	Query       []KeyValue    `json:"query"`
	Tags        []string      `json:"tags"`
	Response    *HttpResponse `json:"response"`
}

type HttpResponse struct {
	Body       string       `json:"body"`
	StatusCode int          `json:"status_code"`
	ID         string       `json:"id"`
	LocalID    int64        `json:"local_id"`
	Headers    []KeyValue   `json:"headers"`
	Tags       []string     `json:"tags"`
	BodySize   int          `json:"body_size"`
	Cookies    []KeyValue   `json:"cookies"`
	Request    *HttpRequest `json:"request"`
}

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func PackageHttpRequest(request *http.Request, proxyID string, reqID int64) (*HttpRequest, error) {
	// replace nil body with empty body
	if request.Body == nil {
		request.Body = io.NopCloser(strings.NewReader(""))
	}

	// create 2 copies of the body - one for the request, one for the raw
	backup, localCopy := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	if _, err := io.Copy(io.MultiWriter(backup, localCopy), request.Body); err != nil {
		return nil, fmt.Errorf("failed to copy request body: %w", err)
	}
	if err := request.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close request body: %w", err)
	}

	// restore original body for entire request write
	request.Body = io.NopCloser(backup)

	resp, _ := PackageHttpResponse(request.Response, proxyID, reqID)

	return &HttpRequest{
		ID:          fmt.Sprintf("%s:%d", proxyID, reqID),
		LocalID:     reqID,
		Method:      request.Method,
		URL:         request.URL.String(),
		Body:        localCopy.String(),
		Host:        request.Host,
		Path:        request.URL.Path,
		QueryString: request.URL.RawQuery,
		Headers:     packageMap(request.Header),
		Query:       packageMap(request.URL.Query()),
		Scheme:      request.URL.Scheme,
		Tags:        tagRequest(request),
		Response:    resp,
	}, nil
}

func tagRequest(req *http.Request) []string {
	tags := []string{} // define as empty array for json friendliness

	if req.Header.Get("Authorization") != "" {
		tags = append(tags, "Auth")
	}

	if req.Header.Get("Cookie") != "" {
		tags = append(tags, "Cookies")
	}

	if tag := tagContentType(req.Header.Get("Content-Type")); tag != "" {
		tags = append(tags, "Request: "+tag)
	}

	return tags
}

func tagResponse(resp *http.Response) []string {
	tags := []string{} // define as empty array for json friendliness

	if resp.Header.Get("Set-Cookie") != "" {
		tags = append(tags, "Set-Cookie")
	}

	if tag := tagContentType(resp.Header.Get("Content-Type")); tag != "" {
		tags = append(tags, "Response: "+tag)
	}

	return tags
}

func tagContentType(ct string) string {
	switch {
	case strings.Contains(ct, "/json"):
		return "JSON"
	case strings.Contains(ct, "/xml"):
		return "XML"
	case strings.Contains(ct, "/html"):
		return "HTML"
	case strings.Contains(ct, "/javascript"):
		return "JS"
	case strings.Contains(ct, "/css"):
		return "CSS"
	case strings.Contains(ct, "/plain"):
		return "Plain"
	default:
		return ""
	}
}

func packageMap(header map[string][]string) []KeyValue {
	headers := []KeyValue{} // nolint
	for key, values := range header {
		for _, value := range values {
			headers = append(headers, KeyValue{
				Key:   key,
				Value: value,
			})
		}
	}
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].Key < headers[j].Key
	})
	return headers
}

func PackageHttpResponse(response *http.Response, proxyID string, reqID int64) (*HttpResponse, error) {

	if response == nil {
		return nil, nil
	}

	// replace nil body with empty body
	if response.Body == nil {
		response.Body = io.NopCloser(strings.NewReader(""))
	}

	// create 2 copies of the body - one for the response, one for the raw
	backup, localCopy := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	if _, err := io.Copy(io.MultiWriter(backup, localCopy), response.Body); err != nil {
		return nil, fmt.Errorf("failed to copy response body: %w", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	// restore original body for entire response write
	response.Body = io.NopCloser(backup)

	cookies := make([]KeyValue, 0)
	for _, cookie := range response.Cookies() {
		cookies = append(cookies, KeyValue{
			Key:   cookie.Name,
			Value: cookie.Value,
		})
	}

	return &HttpResponse{
		ID:         fmt.Sprintf("%s:%d", proxyID, reqID),
		LocalID:    reqID,
		BodySize:   localCopy.Len(),
		Body:       localCopy.String(),
		StatusCode: response.StatusCode,
		Headers:    packageMap(response.Header),
		Tags:       tagResponse(response),
		Cookies:    cookies,
	}, nil
}

func UnpackageHttpRequest(h *HttpRequest) (*http.Request, error) {
	req, err := http.NewRequest(h.Method, h.URL, strings.NewReader(h.Body))
	if err != nil {
		return nil, err
	}
	for _, header := range h.Headers {
		req.Header.Add(header.Key, header.Value)
	}
	return req, nil
}

func UnpackageHttpResponse(h *HttpResponse, req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("not implemented")
}
