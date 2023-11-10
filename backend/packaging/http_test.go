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
				ID:          expectedID,
				LocalID:     requestID,
				Headers:     []KeyValue{},
				Query: []KeyValue{
					{
						Key:   "x",
						Value: "1",
					},
				},
				Tags: []string{},
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
				ID:          expectedID,
				LocalID:     requestID,
				Headers: []KeyValue{
					{
						Key:   "X-Test",
						Value: "test",
					},
				},
				Query: []KeyValue{
					{
						Key:   "x",
						Value: "1",
					},
				},
				Tags: []string{},
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
				ID:          expectedID,
				LocalID:     requestID,
				Headers:     []KeyValue{},
				Query: []KeyValue{
					{
						Key:   "x",
						Value: "1",
					},
				},
				Body: "msg=hello",
				Tags: []string{},
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
				ID:         expectedID,
				LocalID:    requestID,
				Headers:    []KeyValue{},
				Cookies:    []KeyValue{},
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
				ID:      expectedID,
				LocalID: requestID,
				Headers: []KeyValue{
					{
						Key:   "X-Test",
						Value: "test",
					},
				},
				Cookies:    []KeyValue{},
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
				ID:         expectedID,
				LocalID:    requestID,
				Headers:    []KeyValue{},
				Cookies:    []KeyValue{},
				StatusCode: 200,
				Tags:       []string{},
				BodySize:   9,
				Body:       "msg=hello",
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
		})
	}
}
