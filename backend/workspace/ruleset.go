package workspace

import (
	"encoding/json"
	"net/http"
)

type Ruleset []Rule

func (r Ruleset) Match(request *http.Request) bool {
	for _, rule := range r {
		if rule.Match(request) {
			return true
		}
	}
	return false
}

func (r Ruleset) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]Rule(r))
}
