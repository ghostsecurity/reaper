package workspace

import "github.com/ghostsecurity/reaper/backend/packaging"

type Collection struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Name     string    `json:"name"`
	Requests []Request `json:"requests"`
}

type Request struct {
	Name       string                `json:"name"`
	Inner      packaging.HttpRequest `json:"inner"`
	PreScript  string                `json:"pre_script"`
	PostScript string                `json:"post_script"`
}
