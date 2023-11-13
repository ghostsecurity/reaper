package workflow

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/ghostsecurity/reaper/backend/workflow/node"
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
	Name    string
	Status  NodeStatus
	Message string
}

type UpdateM struct {
	Node    string `json:"node"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (u Update) Pack() UpdateM {
	return UpdateM{
		Node:    u.Node.String(),
		Status:  string(u.Status),
		Message: u.Message,
	}
}

var ChildNodeError = errors.New("child node error")

type NodeStatus string

const (
	NodeStatusPending      NodeStatus = "pending"
	NodeStatusRunning      NodeStatus = "running"
	NodeStatusSuccess      NodeStatus = "success"
	NodeStatusError        NodeStatus = "error"
	NodeStatusAborted      NodeStatus = "aborted"
	NodeStatusDisconnected NodeStatus = "disconnected"
)

func (s NodeStatus) IsFinal() bool {
	return s == NodeStatusSuccess || s == NodeStatusError || s == NodeStatusAborted || s == NodeStatusDisconnected
}

func (r *runner) Run(updateChan chan<- Update, output chan<- node.Output) error {

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

	for _, n := range r.workflow.Nodes {
		injections := n.GetInjections()
		if len(injections) == 0 {
			continue
		}
		for _, link := range r.workflow.Links {
			if link.From.Node == n.ID() {
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

	bus := NewBus(startNode, updateChan)

	for _, node := range r.workflow.Nodes {
		if err := bus.AddNode(node); err != nil {
			return fmt.Errorf("failed to add node to bus: %s", err)
		}
	}

	for _, link := range r.workflow.Links {
		if err := bus.AddLink(link); err != nil {
			return fmt.Errorf("failed to add link to bus: %s", err)
		}
	}

	return bus.Run(r.ctx, output)
}
