package workspace

import "net/http"

type Ruleset []Rule

func (r Ruleset) Match(request *http.Request) bool {
	for _, rule := range r {
		if rule.Match(request) {
			return true
		}
	}
	return false
}
