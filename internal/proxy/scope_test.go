package proxy

import "testing"

func TestScope(t *testing.T) {
	scope := NewScope(
		[]string{"acme.com", "example.org"},
		[]string{"special.host.com"},
	)

	tests := []struct {
		host    string
		inScope bool
	}{
		{"acme.com", true},
		{"api.acme.com", true},
		{"deep.api.acme.com", true},
		{"notacme.com", false},
		{"example.org", true},
		{"sub.example.org", true},
		{"example.com", false},
		{"special.host.com", true},
		{"special.host.com:443", true},
		{"other.host.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			if got := scope.InScope(tt.host); got != tt.inScope {
				t.Errorf("InScope(%q) = %v, want %v", tt.host, got, tt.inScope)
			}
		})
	}
}

func TestScopeEmpty(t *testing.T) {
	scope := NewScope(nil, nil)
	if scope.InScope("anything.com") {
		t.Error("empty scope should match nothing")
	}
}
