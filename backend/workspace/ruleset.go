package workspace

import (
	"encoding/json"
	"net/http"
)

type RuleSet []Rule

func (r RuleSet) Match(request *http.Request) bool {
	for _, rule := range r {
		if rule.Match(request) {
			return true
		}
	}
	return false
}

func (r RuleSet) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]Rule(r))
}
