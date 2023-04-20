package app

import (
	"github.com/ghostsecurity/reaper/backend/workspace"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) CreateWorkspace(template *workspace.Workspace) *workspace.Workspace {
	ws := workspace.New()
	ws.Name = template.Name
	ws.Scope = template.Scope
	if err := ws.Save(); err != nil {
		a.logger.Errorf("Failed to write workspace: %s", err)
		a.notifyUser("Failed to write workspace: "+err.Error(), runtime.WarningDialog)
		return nil
	}
	return ws
}

func (a *App) GetWorkspaces() []*workspace.Workspace {
	list, err := workspace.List(a.logger)
	if err != nil {
		a.logger.Errorf("Failed to list workspaces: %s", err)
		return []*workspace.Workspace{}
	}
	return list
}

func (a *App) SaveWorkspace(ws *workspace.Workspace) {
	a.SetWorkspace(ws)
	if err := ws.Save(); err != nil {
		a.logger.Errorf("Failed to save workspace: %s", err)
	}
}

func (a *App) LoadWorkspace(id string) *workspace.Workspace {
	ws, err := workspace.Load(id)
	if err != nil {
		a.logger.Errorf("Failed to load workspace: %s", err)
		return nil
	}
	return ws
}

func (a *App) DeleteWorkspace(id string) {
	if err := workspace.Delete(id); err != nil {
		a.logger.Errorf("Failed to delete workspace: %s", err)
	}
}

func (a *App) SetWorkspace(workspace *workspace.Workspace) {
	a.workspaceMu.Lock()
	defer a.workspaceMu.Unlock()
	a.workspace = workspace
	a.interceptor.SetScope(workspace.InterceptionScope)
}
