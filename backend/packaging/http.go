// Package packaging provides functions to convert to and from http.Request/Response and json marshallable equivalents.
package packaging

import (
	"bufio"
	"bytes"
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
	Request    HttpRequest
	StatusCode int
	ID         int64
	Headers    map[string][]string
}

func PackageHttpRequest(request *http.Request, id int64) HttpRequest {
	var buf bytes.Buffer
	_ = request.Write(&buf)
	return HttpRequest{
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
	}
}

func packageHeaders(header http.Header) map[string][]string {
	headers := make(map[string][]string)
	for k, v := range header {
		headers[k] = v
	}
	return headers
}

func PackageHttpResponse(response *http.Response, id int64) HttpResponse {
	var buf bytes.Buffer
	body, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()
	response.Body = io.NopCloser(bytes.NewBuffer(body))
	_ = response.Write(&buf)
	_ = response.Body.Close()
	response.Body = io.NopCloser(bytes.NewBuffer(body))
	return HttpResponse{
		ID:         id,
		Raw:        buf.String(),
		StatusCode: response.StatusCode,
		Request:    PackageHttpRequest(response.Request, id),
		Headers:    packageHeaders(response.Header),
	}
}

func UnpackageHttpRequest(h *HttpRequest) (*http.Request, error) {
	return http.ReadRequest(bufio.NewReader(strings.NewReader(h.Raw)))
}

func UnpackageHttpResponse(h *HttpResponse, req *http.Request) (*http.Response, error) {
	return http.ReadResponse(bufio.NewReader(strings.NewReader(h.Raw)), req)
}
