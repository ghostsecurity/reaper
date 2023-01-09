package workspace

import (
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestRuleset_Match(t *testing.T) {
	tests := []struct {
		name  string
		rules []Rule
		url   string
		want  bool
	}{
		{
			name:  "empty",
			rules: []Rule{},
			url:   "https://example.com",
			want:  false,
		},
		{
			name: "single match",
			rules: []Rule{
				{
					HostRegex: regexp.MustCompile(`^example\.com`),
				},
			},
			url:  "https://example.com",
			want: true,
		},
		{
			name: "single mismatch",
			rules: []Rule{
				{
					HostRegex: regexp.MustCompile(`^example\.com`),
				},
			},
			url:  "https://whatever.com",
			want: false,
		},
		{
			name: "many, only one need match",
			rules: []Rule{
				{
					HostRegex: regexp.MustCompile(`^whatever\.com`),
				},
				{
					HostRegex: regexp.MustCompile(`^example\.com`),
				},
			},
			url:  "https://example.com",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RuleSet(tt.rules)
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("invalid test url, error = %v", err)
			}
			if got := r.Match(&http.Request{
				URL: u,
			}); got != tt.want {
				t.Errorf("Ruleset.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
