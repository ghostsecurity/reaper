package api

import (
	"github.com/ghostsecurity/reaper/backend/workspace"
)

func (a *API) CreateWorkspace(template *workspace.Workspace) *workspace.Workspace {
	ws := workspace.New()
	ws.Name = template.Name
	ws.Scope = template.Scope
	if err := ws.Save(); err != nil {
		a.logger.Errorf("Failed to write workspace: %s", err)
		a.notify("Failed to write workspace: %s", err)
		return nil
	}
	return ws
}

func (a *API) GetWorkspaces() []*workspace.Workspace {
	list, err := workspace.List(a.logger)
	if err != nil {
		a.logger.Errorf("Failed to list workspaces: %s", err)
		return []*workspace.Workspace{}
	}
	return list
}

func (a *API) SaveWorkspace(ws *workspace.Workspace) {
	a.SetWorkspace(ws)
	if err := ws.Save(); err != nil {
		a.logger.Errorf("Failed to save workspace: %s", err)
	}
}

func (a *API) LoadWorkspace(id string) *workspace.Workspace {
	ws, err := workspace.Load(id)
	if err != nil {
		a.logger.Errorf("Failed to load workspace: %s", err)
		return nil
	}
	return ws
}

func (a *API) DeleteWorkspace(id string) {
	if err := workspace.Delete(id); err != nil {
		a.logger.Errorf("Failed to delete workspace: %s", err)
	}
}

func (a *API) SetWorkspace(workspace *workspace.Workspace) {
	a.workspaceMu.Lock()
	defer a.workspaceMu.Unlock()
	a.workspace = workspace
	a.interceptor.SetScope(workspace.InterceptionScope)
}

func (a *API) GetWorkspace() *workspace.Workspace {
	a.workspaceMu.RLock()
	defer a.workspaceMu.RUnlock()
	return a.workspace
}
