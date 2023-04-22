package workflow

import (
	"fmt"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Workflow struct {
	ID      uuid.UUID             `json:"id"`
	Name    string                `json:"name"`
	Request packaging.HttpRequest `json:"request"`
	Input   Link                  `json:"input"`
	Output  Node                  `json:"output"`
	Error   Node                  `json:"error"`
	Nodes   []Node                `json:"nodes"`
	Links   []Link                `json:"links"`
}

func (w *Workflow) Run(ctx context.Context, updateChan chan<- Update) error {
	if w.Input.To.Node == uuid.Nil {
		return fmt.Errorf("workflow has no input connected")
	}
	return newRunner(ctx, w).Run(updateChan)
}

func (w *Workflow) Validate() error {
	for _, node := range w.Nodes {
		if err := w.validateNode(node); err != nil {
			return err
		}
	}
	if err := w.validateNode(w.Output); err != nil {
		return err
	}
	if err := w.validateNode(w.Error); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) validateNode(node Node) error {
	for _, input := range node.GetOutputs() {
		if err := input.Type.validate(); err != nil {
			return err
		}
	}
	if err := node.GetInput().Type.validate(); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) FindNode(id uuid.UUID) (Node, error) {
	for _, n := range w.Nodes {
		if n.ID() == id {
			return n, nil
		}
	}
	if w.Output.ID() == id {
		return w.Output, nil
	}
	if w.Error.ID() == id {
		return w.Error, nil
	}
	return nil, fmt.Errorf("node not found")
}
