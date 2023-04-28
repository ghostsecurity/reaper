package workflow

import (
	"fmt"
	"io"

	"github.com/ghostsecurity/reaper/backend/workflow/node"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Workflow struct {
	ID          uuid.UUID
	Name        string
	Nodes       []node.Node
	Links       []node.Link
	Positioning map[uuid.UUID]Position
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewWorkflow() *Workflow {
	start := node.NewStart()
	return &Workflow{
		ID:   uuid.New(),
		Name: "New Workflow",
		Nodes: []node.Node{
			start,
		},
		Links: []node.Link{},
		Positioning: map[uuid.UUID]Position{
			start.ID(): Position{X: 50, Y: 50},
		},
	}
}

func (w *Workflow) Run(ctx context.Context, updateChan chan<- Update, stdout, stderr io.Writer) error {
	return newRunner(ctx, w).Run(updateChan, stdout, stderr)
}

func (w *Workflow) Validate() error {
	for _, node := range w.Nodes {
		if err := node.Validate(nil); err != nil {
			return err
		}
	}
	for _, link := range w.Links {
		var foundFrom bool
		var foundTo bool
		for _, node := range w.Nodes {
			if node.ID() == link.From.Node {
				foundFrom = true
				if _, ok := node.GetOutputs().FindByName(link.From.Connector); !ok {
					return fmt.Errorf("from node %s has no output '%s'", link.From.Node, link.From.Connector)
				}
			}
			if node.ID() == link.To.Node {
				foundTo = true
				if _, ok := node.GetInputs().FindByName(link.To.Connector); !ok {
					return fmt.Errorf("to node %s has no input '%s'", link.To.Node, link.To.Connector)
				}
			}
		}
		if !foundFrom {
			return fmt.Errorf("link: from node %s not found", link.From.Node)
		}
		if !foundTo {
			return fmt.Errorf("link: to node %s not found", link.To.Node)
		}
		if link.From.Node == link.To.Node {
			return fmt.Errorf("link: from and to node are the same")
		}
	}
	return nil
}

func (w *Workflow) FindNode(id uuid.UUID) (node.Node, error) {
	for _, n := range w.Nodes {
		if n.ID() == id {
			return n, nil
		}
	}
	return nil, fmt.Errorf("node not found")
}
