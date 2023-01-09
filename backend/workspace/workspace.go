package workspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/kirsle/configdir"

	"github.com/google/uuid"
)

const (
	configDirName     = "reaper"
	workspacesDirName = "workspaces"
)

func getDir() (string, error) {
	wsDir := configdir.LocalConfig(configDirName, workspacesDirName)
	if err := configdir.MakePath(wsDir); err != nil {
		return "", fmt.Errorf("failed to create workspace dir: %w", err)
	}
	return wsDir, nil
}

type Workspace struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Scope      Scope      `json:"scope"`
	Collection Collection `json:"collection"`
	Tree       Tree       `json:"tree"`
	// TODO: flows
}

func New() *Workspace {
	return &Workspace{
		ID:   uuid.New().String(),
		Name: "Untitled Workspace",
		Scope: Scope{
			Include: Ruleset{
				{
					HostRegex: regexp.MustCompile(`^(api\.)?ghostbank\.net$`), // TODO: remove this
				},
			},
		},
	}
}

func List() ([]*Workspace, error) {
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
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".workspace") {
			continue
		}
		workspace, err := loadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to load workspace %s: %w", entry.Name(), err)
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
	return filepath.Join(dir, id+".workspace"), nil
}

func Load(id string) (*Workspace, error) {

	file, err := getWorkspacePath(id)
	if err != nil {
		return nil, fmt.Errorf("failed to locate workspace path: %w", err)
	}

	return loadFile(file)
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

	return &workspace, nil
}

func (w *Workspace) Save() error {

	file, err := getWorkspacePath(w.ID)
	if err != nil {
		return fmt.Errorf("failed to locate workspace path: %w", err)
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

func (w *Workspace) UpdateTree(request *http.Request) *Tree {
	w.Tree.Update(request)
	return &w.Tree
}
