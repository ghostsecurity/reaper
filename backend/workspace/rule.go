package workspace

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

type Rule struct {
	ID        int            `json:"id"`
	Protocol  string         `json:"protocol"` // e.g. http, https
	HostRegex *regexp.Regexp `json:"host"`
	PathRegex *regexp.Regexp `json:"path"`
	Ports     PortList       `json:"ports"`
}

func (r Rule) MarshalJSON() ([]byte, error) {
	var hostPattern string
	if r.HostRegex != nil {
		hostPattern = r.HostRegex.String()
	}
	var pathPattern string
	if r.PathRegex != nil {
		pathPattern = r.PathRegex.String()
	}
	if r.Ports == nil {
		r.Ports = PortList{}
	}
	return json.Marshal(map[string]interface{}{
		"protocol": r.Protocol,
		"host":     hostPattern,
		"path":     pathPattern,
		"ports":    r.Ports,
	})
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Protocol = raw["protocol"].(string)
	if pattern, ok := raw["host"].(string); ok && pattern != "" {
		r.HostRegex = regexp.MustCompile(pattern)
	}
	if pattern, ok := raw["path"].(string); ok && pattern != "" {
		r.PathRegex = regexp.MustCompile(pattern)
	}
	if ports, ok := raw["ports"].([]interface{}); ok {
		r.Ports = make(PortList, len(ports))
		for i, port := range ports {
			r.Ports[i] = int(port.(float64))
		}
	}
	return nil
}

func (r *Rule) Match(request *http.Request) bool {
	if request == nil {
		return false
	}
	if r.HostRegex != nil {
		if !r.HostRegex.MatchString(request.URL.Hostname()) {
			return false
		}
	}
	if r.PathRegex != nil {
		if !r.PathRegex.MatchString(request.URL.Path) {
			return false
		}
	}

	if r.Protocol != "" && r.Protocol != request.URL.Scheme {
		return false
	}

	if len(r.Ports) > 0 {
		portStr := request.URL.Port()
		var port int
		if portStr != "" {
			var err error
			port, err = strconv.Atoi(request.URL.Port())
			if err != nil {
				return false
			}
		} else {
			switch request.URL.Scheme {
			case "http":
				port = 80
			case "https":
				port = 443
			}
		}
		if !r.Ports.Match(port) {
			return false
		}
	}

	return true
}
