package workspace

import (
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestScope_Includes(t *testing.T) {
	tests := []struct {
		name     string
		includes Ruleset
		excludes Ruleset
		url      string
		want     bool
	}{
		{
			name: "empty",
			url:  "https://example.com",
			want: true,
		},
		{
			name: "include match",
			includes: Ruleset{
				{
					HostRegex: regexp.MustCompile(`^example\.com`),
				},
			},
			url:  "https://example.com",
			want: true,
		},
		{
			name: "include match, exclude mismatch",
			includes: Ruleset{
				{
					HostRegex: regexp.MustCompile(`^example\.com`),
				},
			},
			excludes: Ruleset{
				{
					HostRegex: regexp.MustCompile(`^whatever\.com`),
				},
			},
			url:  "https://example.com",
			want: true,
		},
		{
			name: "include match, exclude match",
			includes: Ruleset{
				{
					HostRegex: regexp.MustCompile(`^example\.com`),
				},
			},
			excludes: Ruleset{
				{
					HostRegex: regexp.MustCompile(`^example\.com`),
				},
			},
			url:  "https://example.com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := Scope{
				Include: tt.includes,
				Exclude: tt.excludes,
			}
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("invalid test url, error = %v", err)
			}
			if got := scope.Includes(&http.Request{
				URL: u,
			}); got != tt.want {
				t.Errorf("Scope.Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}
