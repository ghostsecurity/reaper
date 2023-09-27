package api

import "github.com/ghostsecurity/reaper/backend/workspace"

type API struct{}

func New() *API {
	return &API{}
}

type Something struct {
	Msg string `json:"msg"`
}

func (a *API) Test() (workspace.Workspace, error) {
	return workspace.Workspace{}, nil
}

func (a *API) HelloWorld() (string, error) {
	return "Hello world", nil
}

func (a *API) List() []string {
	return []string{"a", "b", "c"}
}
