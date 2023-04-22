package workflow

import (
	"fmt"

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
)

func (r *runner) Run(updateChan chan<- Update) error {

	if r.workflow == nil {
		return fmt.Errorf("workflow is nil")
	}

	if r.workflow.Input.To.Node == uuid.Nil {
		return fmt.Errorf("workflow has no input connected")
	}

	start, err := r.workflow.FindNode(r.workflow.Input.To.Node)
	if err != nil {
		return err
	}

	for _, node := range r.workflow.Nodes {
		updateChan <- Update{
			Node:    node.ID(),
			Status:  NodeStatusPending,
			Message: "Waiting for input(s)...",
		}
	}

	return r.RunNode(start, map[uuid.UUID]Transmission{
		r.workflow.Input.To.Connector: NewRequestTransmission(r.workflow.Request),
	}, updateChan, true)
}

func (r *runner) RunNode(node Node, params map[uuid.UUID]Transmission, updateChan chan<- Update, lastInput bool) error {

	updateChan <- Update{
		Node:    node.ID(),
		Status:  NodeStatusRunning,
		Message: "In Progress...",
	}

	outputChan, errChan := node.Run(r.ctx, params)
	for {
		select {
		case <-r.ctx.Done():
			updateChan <- Update{
				Node:    node.ID(),
				Status:  NodeStatusError,
				Message: fmt.Sprintf("Workflow cancelled: %s", r.ctx.Err()),
			}
			return r.ctx.Err()
		case output, ok := <-outputChan:
			if !ok {
				if lastInput {
					updateChan <- Update{
						Node:    node.ID(),
						Status:  NodeStatusSuccess,
						Message: "Operation complete.",
					}
				}
				return nil
			}
			if lastInput && output.Complete {
				updateChan <- Update{
					Node:    node.ID(),
					Status:  NodeStatusSuccess,
					Message: "Operation complete.",
				}
			}
			for _, link := range r.workflow.Links {
				if link.From.Node == node.ID() && link.From.Connector == output.OutputID {
					nextNode, err := r.workflow.FindNode(link.To.Node)
					if err != nil {
						return err
					}
					if err := r.RunNode(nextNode, map[uuid.UUID]Transmission{
						link.To.Connector: output.Data,
					}, updateChan, output.Complete); err != nil {
						return fmt.Errorf("error with node '%s': %s", nextNode.Name(), err)
					}
				}
			}
		case err, ok := <-errChan:
			if ok {
				updateChan <- Update{
					Node:    node.ID(),
					Status:  NodeStatusError,
					Message: fmt.Sprintf("Error: %s", err),
				}
				return err
			}
		}
	}
}
