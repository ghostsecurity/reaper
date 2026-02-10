package proxy

import "strings"

// Scope determines whether a host is in scope for MITM interception.
type Scope struct {
	domains []string // suffix match (e.g. "acme.com" matches "api.acme.com" and "acme.com")
	hosts   []string // exact match
}

func NewScope(domains, hosts []string) *Scope {
	return &Scope{
		domains: domains,
		hosts:   hosts,
	}
}

func (s *Scope) InScope(host string) bool {
	// Strip port if present
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	host = strings.ToLower(host)

	for _, h := range s.hosts {
		if strings.ToLower(h) == host {
			return true
		}
	}

	for _, d := range s.domains {
		d = strings.ToLower(d)
		if host == d || strings.HasSuffix(host, "."+d) {
			return true
		}
	}

	return false
}
