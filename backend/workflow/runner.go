package workflow

import (
	"fmt"
	"io"

	"github.com/ghostsecurity/reaper/backend/workflow/node"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type runner struct {
	workflow *Workflow
	ctx      context.Context
}

func newRunner(ctx context.Context, w *Workflow) *runner {
	return &runner{
		workflow: w,
		ctx:      ctx,
	}
}

type Update struct {
	Node    uuid.UUID
	Status  NodeStatus
	Message string
}

type NodeStatus string

const (
	NodeStatusPending NodeStatus = "pending"
	NodeStatusRunning NodeStatus = "running"
	NodeStatusSuccess NodeStatus = "success"
	NodeStatusError   NodeStatus = "error"
	NodeStatusAborted NodeStatus = "aborted"
)

func (r *runner) Run(updateChan chan<- Update, stdout, stderr io.Writer) error {

	if r.workflow == nil {
		return fmt.Errorf("workflow is nil")
	}

	var startNode node.Node
	for _, n := range r.workflow.Nodes {
		if n.Type() == node.TypeStart {
			startNode = n
			break
		}
	}

	if startNode == nil {
		return fmt.Errorf("workflow has no start node")
	}

	for _, node := range r.workflow.Nodes {
		updateChan <- Update{
			Node:    node.ID(),
			Status:  NodeStatusPending,
			Message: "Waiting for input(s)...",
		}
	}

	for _, node := range r.workflow.Nodes {
		injections := node.GetInjections()
		if len(injections) == 0 {
			continue
		}
		for _, link := range r.workflow.Links {
			if link.From.Node == node.ID() {
				for name, trans := range injections {
					if name == link.From.Connector {
						target, err := r.workflow.FindNode(link.To.Node)
						if err != nil {
							break
						}
						if _, ok := target.GetInputs().FindByName(link.To.Connector); !ok {
							break
						}
						if err := target.AddStaticInputValue(link.To.Connector, trans); err != nil {
							return fmt.Errorf("failed to inject value: %s", err)
						}
					}
				}
			}
		}
	}

	return r.RunNode(startNode, nil, updateChan, true, stdout, stderr)
}

func (r *runner) RunNode(n node.Node, params map[string]transmission.Transmission, updateChan chan<- Update, lastInput bool, stdout, stderr io.Writer) error {

	if err := n.Validate(params); err != nil {
		return fmt.Errorf("invalid node: %s", err)
	}

	updateChan <- Update{
		Node:    n.ID(),
		Status:  NodeStatusRunning,
		Message: "In Progress...",
	}

	outputChan, errChan := n.Run(r.ctx, params, stdout, stderr)
	defer waitForChannels(outputChan, errChan)
	for {
		select {
		case <-r.ctx.Done():
			updateChan <- Update{
				Node:    n.ID(),
				Status:  NodeStatusError,
				Message: fmt.Sprintf("Workflow cancelled: %s", r.ctx.Err()),
			}
			return r.ctx.Err()
		case output, ok := <-outputChan:
			if !ok {
				if lastInput {
					updateChan <- Update{
						Node:    n.ID(),
						Status:  NodeStatusSuccess,
						Message: "Operation complete.",
					}
				}
				return nil
			}
			if lastInput && output.Complete {
				updateChan <- Update{
					Node:    n.ID(),
					Status:  NodeStatusSuccess,
					Message: "Operation complete.",
				}
			}
			for _, link := range r.workflow.Links {
				if link.From.Node == n.ID() && link.From.Connector == output.OutputName {
					nextNode, err := r.workflow.FindNode(link.To.Node)
					if err != nil {
						return err
					}
					if err := r.RunNode(nextNode, map[string]transmission.Transmission{
						link.To.Connector: output.Data,
					}, updateChan, output.Complete, stdout, stderr); err != nil {
						return fmt.Errorf("error with node '%s': %s", nextNode.Name(), err)
					}
				}
			}
		case err, ok := <-errChan:
			if ok {
				updateChan <- Update{
					Node:    n.ID(),
					Status:  NodeStatusError,
					Message: fmt.Sprintf("Error: %s", err),
				}
				return err
			}
		}
	}
}

func waitForChannels(c <-chan node.OutputInstance, errChan <-chan error) {
	for range c {
	}
	for range errChan {
	}
}
