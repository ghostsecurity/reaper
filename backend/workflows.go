package backend

import (
	"bytes"
	"context"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/google/uuid"

	"github.com/ghostsecurity/reaper/backend/workflow"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
)

func (a *App) RunWorkflow(w *workflow.WorkflowM) {
	if a.runningWorkflowID != uuid.Nil {
		a.Error("Workflow already running", "There is already a workflow running. Please cancel it or wait for it to finish.")
		return
	}
	flow, err := w.Unpack()
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	a.workflowContextCancel = cancel
	stdout := os.Stdout // bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	updateChan := make(chan workflow.Update)
	go func() {
		for x := range updateChan {
			_ = x
		}
	}()
	runtime.EventsEmit(a.ctx, EventWorkflowStarted, w.ID)
	defer runtime.EventsEmit(a.ctx, EventWorkflowFinished, w.ID)
	if err := flow.Run(ctx, updateChan, stdout, stderr); err != nil {
		a.Error("Workflow error", err.Error())
		return
	}
}

func (a *App) StopWorkflow(w *workflow.WorkflowM) {
	if a.workflowContextCancel != nil {
		a.workflowContextCancel()
	}
}

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
