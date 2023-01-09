package workspace

import (
	"net/http"
)

type Scope struct {
	// Include is a list of rules, which, if matched, will result in a request being included in the scope
	// An empty list will result in all requests being included, unless the Exclude list is used to exclude requests
	// ANY rule being matched will result in a request being included
	Include RuleSet `json:"include"`
	// the Exclude ruleset is used to exclude items that have previously been included
	// ANY rule being matched will result in a request being excluded
	Exclude RuleSet `json:"exclude"`
}

func (s Scope) Includes(request *http.Request) bool {
	return (len(s.Include) == 0 || s.Include.Match(request)) && !s.Exclude.Match(request)
}
