package packaging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestPackaging(t *testing.T) {

	proxyID := uuid.New().String()
	requestID := int64(123)
	expectedID := fmt.Sprintf("%s:%d", proxyID, requestID)

	tests := []struct {
		name  string
		input func() *http.Request
		want  HttpRequest
	}{
		{
			name: "get",
			input: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com/blah?x=1", nil)
				return req
			},
			want: HttpRequest{
				Method:      "GET",
				URL:         "https://example.com/blah?x=1",
				Host:        "example.com",
				Path:        "/blah",
				QueryString: "x=1",
				Scheme:      "https",
				Raw:         "GET /blah?x=1 HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\n\r\n",
				ID:          expectedID,
				LocalID:     requestID,
				Headers:     map[string][]string{},
				Query:       map[string][]string{"x": {"1"}},
				Tags:        []string{},
			},
		},
		{
			name: "get with extra header",
			input: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com/blah?x=1", nil)
				req.Header.Set("X-Test", "test")
				return req
			},
			want: HttpRequest{
				Method:      "GET",
				URL:         "https://example.com/blah?x=1",
				Host:        "example.com",
				Path:        "/blah",
				QueryString: "x=1",
				Scheme:      "https",
				Raw:         "GET /blah?x=1 HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nX-Test: test\r\n\r\n",
				ID:          expectedID,
				LocalID:     requestID,
				Headers:     map[string][]string{"X-Test": {"test"}},
				Query:       map[string][]string{"x": {"1"}},
				Tags:        []string{},
			},
		},
		{
			name: "post with body",
			input: func() *http.Request {
				req, _ := http.NewRequest("POST", "https://example.com/blah?x=1", strings.NewReader("msg=hello"))
				return req
			},
			want: HttpRequest{
				Method:      "POST",
				URL:         "https://example.com/blah?x=1",
				Host:        "example.com",
				Path:        "/blah",
				QueryString: "x=1",
				Scheme:      "https",
				Raw:         "POST /blah?x=1 HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 9\r\n\r\nmsg=hello",
				ID:          expectedID,
				LocalID:     requestID,
				Headers:     map[string][]string{},
				Query:       map[string][]string{"x": {"1"}},
				Tags:        []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input()
			got, err := PackageHttpRequest(input, proxyID, requestID)
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, *got)
			var buf bytes.Buffer
			require.NoError(t, input.Write(&buf))
			assert.Equal(t, tt.want.Raw, buf.String())
		})
	}
}

func TestResponsePackaging(t *testing.T) {

	proxyID := uuid.New().String()
	requestID := int64(123)
	expectedID := fmt.Sprintf("%s:%d", proxyID, requestID)

	tests := []struct {
		name  string
		input *http.Response
		want  HttpResponse
	}{
		{
			name: "get",
			input: &http.Response{
				StatusCode: 200,
				Header:     http.Header{},
				Body:       nil,
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			want: HttpResponse{
				Raw:        "HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n",
				ID:         expectedID,
				LocalID:    requestID,
				Headers:    map[string][]string{},
				StatusCode: 200,
				Tags:       []string{},
			},
		},
		{
			name: "get with headers",
			input: &http.Response{
				StatusCode: 200,
				Header: http.Header{
					"X-Test": {"test"},
				},
				Body:       nil,
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			want: HttpResponse{
				Raw:     "HTTP/1.1 200 OK\r\nX-Test: test\r\nContent-Length: 0\r\n\r\n",
				ID:      expectedID,
				LocalID: requestID,
				Headers: map[string][]string{
					"X-Test": {"test"},
				},
				StatusCode: 200,
				Tags:       []string{},
			},
		},
		{
			name: "post with body",
			input: &http.Response{
				StatusCode: 200,
				Header:     http.Header{},
				Body:       io.NopCloser(strings.NewReader("msg=hello")),
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			want: HttpResponse{
				Raw:        "HTTP/1.1 200 OK\r\nConnection: close\r\n\r\nmsg=hello",
				ID:         expectedID,
				LocalID:    requestID,
				Headers:    map[string][]string{},
				StatusCode: 200,
				Tags:       []string{},
				BodySize:   9,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PackageHttpResponse(tt.input, proxyID, requestID)
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, *got)
			var buf bytes.Buffer
			require.NoError(t, tt.input.Write(&buf))
			assert.Equal(t, tt.want.Raw, buf.String())
		})
	}
}
