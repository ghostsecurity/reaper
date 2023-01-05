package workspace

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"regexp"
)

type Workspace struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Scope      Scope      `json:"scope"`
	Collection Collection `json:"collection"`
	Tree       Tree       `json:"tree"`
	// TODO: flows
}

func New() *Workspace {
	return &Workspace{
		ID:   uuid.New(),
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
	// TODO
	return nil, fmt.Errorf("not implemented")
}

func LoadPrevious() (*Workspace, error) {

	// TODO: load previous workspace from disk

	return New(), nil
}

func Load(id uuid.UUID) (*Workspace, error) {
	// TODO
	return nil, fmt.Errorf("not implemented")
}

func (w *Workspace) Save() error {
	// TODO save workspace to disk as json
	return fmt.Errorf("not implemented")
}

func (w *Workspace) UpdateTree(request *http.Request) *Tree {
	w.Tree.Update(request)
	return &w.Tree
}
