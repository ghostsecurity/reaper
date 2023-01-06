// Package packaging provides functions to convert to and from http.Request/Response and json marshallable equivalents.
package packaging

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpRequest struct {
	Method      string
	URL         string
	Host        string
	Path        string
	QueryString string
	Scheme      string
	Raw         string
	ID          int64
	Headers     map[string][]string
	Query       map[string][]string
}

type HttpResponse struct {
	Raw        string
	StatusCode int
	ID         int64
	Headers    map[string][]string
}

func PackageHttpRequest(request *http.Request, id int64) (*HttpRequest, error) {

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

	// write entire request
	var buf bytes.Buffer
	if err := request.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	// restore original body for local copy
	request.Body = io.NopCloser(localCopy)

	return &HttpRequest{
		ID:          id,
		Method:      request.Method,
		URL:         request.URL.String(),
		Raw:         buf.String(),
		Host:        request.Host,
		Path:        request.URL.Path,
		QueryString: request.URL.RawQuery,
		Headers:     packageHeaders(request.Header),
		Query:       request.URL.Query(),
		Scheme:      request.URL.Scheme,
	}, nil
}

func packageHeaders(header http.Header) map[string][]string {
	headers := make(map[string][]string)
	for k, v := range header {
		headers[k] = v
	}
	return headers
}

func PackageHttpResponse(response *http.Response, id int64) (*HttpResponse, error) {

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

	// write entire response
	var buf bytes.Buffer
	if err := response.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write response: %w", err)
	}

	// restore original body for local copy
	response.Body = io.NopCloser(localCopy)

	return &HttpResponse{
		ID:         id,
		Raw:        buf.String(),
		StatusCode: response.StatusCode,
		Headers:    packageHeaders(response.Header),
	}, nil
}

func UnpackageHttpRequest(h *HttpRequest) (*http.Request, error) {
	return http.ReadRequest(bufio.NewReader(strings.NewReader(h.Raw)))
}

func UnpackageHttpResponse(h *HttpResponse, req *http.Request) (*http.Response, error) {
	return http.ReadResponse(bufio.NewReader(strings.NewReader(h.Raw)), req)
}
