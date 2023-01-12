package workspace

import "github.com/ghostsecurity/reaper/backend/packaging"

type Collection struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Requests []Request `json:"requests"`
}

type Request struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	Inner      packaging.HttpRequest `json:"inner"`
	PreScript  string                `json:"pre_script"`
	PostScript string                `json:"post_script"`
}
