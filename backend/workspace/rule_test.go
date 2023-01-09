package workspace

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestRule_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		r       Rule
		want    []byte
		wantErr bool
	}{
		{
			name: "empty",
			r:    Rule{},
			want: []byte(`{"host":"","path":"","ports":[],"protocol":""}`),
		},
		{
			name: "full",
			r: Rule{
				Protocol:     "https:",
				HostRegexRaw: `^example\.com`,
				HostRegex:    regexp.MustCompile(`^example\.com`),
				PathRegexRaw: `/foo`,
				PathRegex:    regexp.MustCompile(`/foo`),
				Ports: PortList{
					80,
					443,
				},
			},
			want: []byte(`{"host":"^example\\.com","path":"/foo","ports":[80,443],"protocol":"https:"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rule.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("Rule.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestRule_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Rule
	}{
		{
			name:  "empty",
			input: `{"host":"","path":"","ports":[],"protocol":""}`,
			want:  Rule{},
		},
		{
			name:  "full",
			input: `{"host":"^example\\.com","path":"/foo","ports":[80,443],"protocol":"https"}`,
			want: Rule{
				Protocol:  "https",
				HostRegex: regexp.MustCompile(`^example\.com`),
				PathRegex: regexp.MustCompile(`/foo`),
				Ports: PortList{
					80,
					443,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Rule
			if err := json.Unmarshal([]byte(tt.input), &got); err != nil {
				t.Errorf("Rule.UnmarshalJSON() error = %v", err)
			}
			if got.Protocol != tt.want.Protocol {
				t.Errorf("Rule.UnmarshalJSON() protocol = %v, want %v", got.Protocol, tt.want.Protocol)
			}
			if len(got.Ports) != len(tt.want.Ports) {
				t.Errorf("Rule.UnmarshalJSON() ports = %v, want %v", got.Ports, tt.want.Ports)
			} else {
				for i, port := range got.Ports {
					if port != tt.want.Ports[i] {
						t.Errorf("Rule.UnmarshalJSON() ports[%d] = %v, want %v", i, port, tt.want.Ports[i])
					}
				}
			}
			if tt.want.HostRegex == nil || got.HostRegex == nil {
				if tt.want.HostRegex != got.HostRegex {
					t.Errorf("Rule.UnmarshalJSON() hostRegex = %v, want %v", got.HostRegex, tt.want.HostRegex)
				}
			} else if got.HostRegex.String() != tt.want.HostRegex.String() {
				t.Errorf("Rule.UnmarshalJSON() hostRegex = %v, want %v", got.HostRegex.String(), tt.want.HostRegex.String())
			}
			if tt.want.PathRegex == nil || got.PathRegex == nil {
				if tt.want.PathRegex != got.PathRegex {
					t.Errorf("Rule.UnmarshalJSON() pathRegex = %v, want %v", got.PathRegex, tt.want.PathRegex)
				}
			} else if got.PathRegex.String() != tt.want.PathRegex.String() {
				t.Errorf("Rule.UnmarshalJSON() pathRegex = %v, want %v", got.PathRegex.String(), tt.want.PathRegex.String())
			}
		})
	}
}

func TestRule_Match(t *testing.T) {

	tests := []struct {
		name string
		r    Rule
		url  string
		want bool
	}{
		{
			name: "empty",
			r:    Rule{},
			url:  "http://example.com",
			want: true,
		},
		{
			name: "protocol with match",
			r: Rule{
				Protocol: "https",
			},
			url:  "https://example.com",
			want: true,
		},
		{
			name: "protocol without match",
			r: Rule{
				Protocol: "http",
			},
			url:  "https://example.com",
			want: false,
		},
		{
			name: "host with match",
			r: Rule{
				HostRegex: regexp.MustCompile(`^example\.com`),
			},
			url:  "https://example.com",
			want: true,
		},
		{
			name: "host with optional subdomain parent match",
			r: Rule{
				HostRegex: regexp.MustCompile(`^(www\.)?example\.com`),
			},
			url:  "https://example.com",
			want: true,
		},
		{
			name: "host with optional subdomain child match",
			r: Rule{
				HostRegex: regexp.MustCompile(`^(www\.)?example\.com`),
			},
			url:  "https://www.example.com",
			want: true,
		},
		{
			name: "host exact mismatch",
			r: Rule{
				HostRegex: regexp.MustCompile(`^example\.com`),
			},
			url:  "https://zexample.com",
			want: false,
		},
		{
			name: "path with match",
			r: Rule{
				PathRegex: regexp.MustCompile(`/foo`),
			},
			url:  "https://example.com/foo",
			want: true,
		},
		{
			name: "path with match",
			r: Rule{
				PathRegex: regexp.MustCompile(`^/foo$`),
			},
			url:  "https://example.com/foo?x=1",
			want: true,
		},
		{
			name: "path without exact match",
			r: Rule{
				PathRegex: regexp.MustCompile(`^/foo$`),
			},
			url:  "https://example.com/foo/bar?x=1",
			want: false,
		},
		{
			name: "port with implied match",
			r: Rule{
				Ports: []int{443},
			},
			url:  "https://example.com",
			want: true,
		},
		{
			name: "port with implied mismatch",
			r: Rule{
				Ports: []int{80},
			},
			url:  "https://example.com",
			want: false,
		},
		{
			name: "port with explicit match",
			r: Rule{
				Ports: []int{443},
			},
			url:  "https://example.com:443",
			want: true,
		},
		{
			name: "port with explicit mismatch",
			r: Rule{
				Ports: []int{443},
			},
			url:  "https://example.com:123",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("invalid test url, error = %v", err)
			}
			if got := tt.r.Match(&http.Request{
				URL: u,
			}); got != tt.want {
				t.Errorf("Rule.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
