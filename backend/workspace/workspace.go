package workspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/ghostsecurity/reaper/backend/workflow"

	"github.com/kirsle/configdir"

	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/google/uuid"
)

const (
	configDirName         = "reaper"
	workspacesDirName     = "workspaces"
	workspaceSettingsFile = "workspace.json"
)

func getDir() (string, error) {
	wsDir := configdir.LocalConfig(configDirName, workspacesDirName)
	if err := configdir.MakePath(wsDir); err != nil {
		return "", fmt.Errorf("failed to create workspace dir: %w", err)
	}
	return wsDir, nil
}

type Workspace struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	Scope             Scope                `json:"scope"`
	InterceptionScope Scope                `json:"interception_scope"`
	Collection        Collection           `json:"collection"`
	Tree              Tree                 `json:"tree"`
	Workflows         []workflow.WorkflowM `json:"workflows"`
	mu                sync.Mutex
}

func New() *Workspace {
	return &Workspace{
		ID:   uuid.New().String(),
		Name: "Untitled Workspace",
	}
}

func List(logger *log.Logger) ([]*Workspace, error) {
	dir, err := getDir()
	if err != nil {
		return nil, fmt.Errorf("failed to access workspace directory: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read workspace directory: %w", err)
	}

	// sort by last modified
	sort.Slice(entries, func(i, j int) bool {
		ii, err := entries[i].Info()
		if err != nil {
			return false
		}
		ji, err := entries[j].Info()
		if err != nil {
			return false
		}
		return ii.ModTime().After(ji.ModTime())
	})

	// we initialise the empty slice here for json encoding friendliness later
	workspaces := []*Workspace{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		workspaceSettings := filepath.Join(dir, entry.Name(), workspaceSettingsFile)

		if _, err := os.Stat(workspaceSettings); os.IsNotExist(err) {
			continue
		}

		workspace, err := loadFile(workspaceSettings)
		if err != nil {
			logger.Errorf("Failed to read workspace '%s': %s", entry.Name(), err)
			continue
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

func getWorkspacePath(id string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("missing workspace id")
	}
	dir, err := getDir()
	if err != nil {
		return "", fmt.Errorf("failed to access workspace directory: %w", err)
	}
	return filepath.Join(dir, id, workspaceSettingsFile), nil
}

func Load(id string) (*Workspace, error) {

	file, err := getWorkspacePath(id)
	if err != nil {
		return nil, fmt.Errorf("failed to locate workspace path: %w", err)
	}

	return loadFile(file)
}

func Delete(id string) error {
	workspacePath := filepath.Join(configdir.LocalConfig(configDirName, workspacesDirName), id)
	return os.RemoveAll(workspacePath)
}

func loadFile(file string) (*Workspace, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open workspace file: %w", err)
	}
	defer func() { _ = f.Close() }()

	var workspace Workspace
	if err := json.NewDecoder(f).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("failed to decode workspace: %w", err)
	}

	if workspace.Workflows == nil {
		workspace.Workflows = []workflow.WorkflowM{}
	}

	return &workspace, nil
}

func (w *Workspace) Save() error {

	w.mu.Lock()
	defer w.mu.Unlock()

	file, err := getWorkspacePath(w.ID)
	if err != nil {
		return fmt.Errorf("failed to locate workspace path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %w", err)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create/truncate workspace file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := json.NewEncoder(f).Encode(w); err != nil {
		return fmt.Errorf("failed to encode workspace: %w", err)
	}

	return nil
}

func (w *Workspace) UpdateTree(request *http.Request) (*Tree, bool) {

	w.mu.Lock()
	defer w.mu.Unlock()

	changed := w.Tree.Update(request)
	return &w.Tree, changed
}
