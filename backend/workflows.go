package backend

import (
	"github.com/ghostsecurity/reaper/backend/workflow"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
)

func (a *App) CreateWorkflow() *workflow.WorkflowM {
	w, err := workflow.NewWorkflow().Pack()
	if err != nil {
		return nil
	}
	return w
}

func (a *App) CreateNode(nodeType int) *workflow.NodeM {
	switch node.Type(nodeType) {
	case node.TypeFuzzer:
		f, err := workflow.ToNodeM(node.NewFuzzer())
		if err != nil {
			return nil
		}
		return f
	case node.TypeOutput:
		o, err := workflow.ToNodeM(node.NewOutput())
		if err != nil {
			return nil
		}
		return o
	case node.TypeStatusFilter:
		s, err := workflow.ToNodeM(node.NewStatusFilter())
		if err != nil {
			return nil
		}
		return s
	case node.TypeRequest:
		s, err := workflow.ToNodeM(node.NewRequest())
		if err != nil {
			return nil
		}
		return s
	case node.TypeSender:
		s, err := workflow.ToNodeM(node.NewSender())
		if err != nil {
			return nil
		}
		return s
	default:
		return nil
	}
}
