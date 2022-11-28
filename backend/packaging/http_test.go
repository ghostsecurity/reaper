package packaging

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestRequestPackaging(t *testing.T) {
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
				Raw:         "GET /blah?x=1 HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\n\r\n",
				ID:          123,
				Headers:     map[string][]string{},
				Query:       map[string][]string{"x": {"1"}},
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
				Raw:         "GET /blah?x=1 HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nX-Test: test\r\n\r\n",
				ID:          123,
				Headers:     map[string][]string{"X-Test": {"test"}},
				Query:       map[string][]string{"x": {"1"}},
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
				Raw:         "POST /blah?x=1 HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 9\r\n\r\nmsg=hello",
				ID:          123,
				Headers:     map[string][]string{},
				Query:       map[string][]string{"x": {"1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PackageHttpRequest(tt.input(), 123)
			assert.Equal(t, tt.want, got)
		})
	}
}
