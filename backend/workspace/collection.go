package workspace

import "github.com/ghostsecurity/reaper/backend/packaging"

type Collection struct {
	Groups []Group
}

type Group struct {
	Name     string
	Requests []Request
}

type Request struct {
	Name       string
	Inner      packaging.HttpRequest
	PreScript  string
	PostScript string
}
