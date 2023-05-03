package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"

	"github.com/ghostsecurity/reaper/backend/packaging"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/google/uuid"

	"github.com/ghostsecurity/reaper/backend/workflow"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
)

func (a *App) IgnoreThisUsedBindings(_ node.OutputM) workflow.UpdateM {
	return workflow.UpdateM{}
}

func (a *App) SelectFile(title string) (string, error) {
	if title == "" {
		title = "Select a file"
	}
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
}

func (a *App) RunWorkflow(w *workflow.WorkflowM) {
	a.workflowMu.Lock()
	defer a.workflowMu.Unlock()
	if a.runningWorkflowID != uuid.Nil {
		a.Error("Workflow already running", "There is already a workflow running. Please cancel it or wait for it to finish.")
		return
	}
	flow, err := w.Unpack()
	if err != nil {
		a.Error("Workflow error", err.Error())
		return
	}
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	a.workflowContextCancel = cancel
	updateChan := make(chan workflow.Update)
	outputChan := make(chan node.Output)
	defer close(updateChan)
	defer close(outputChan)
	go func() {
		for update := range updateChan {
			runtime.EventsEmit(a.ctx, EventWorkflowUpdate, update.Pack())

			if n, err := flow.FindNode(update.Node); err == nil {
				runtime.EventsEmit(a.ctx, EventWorkflowOutput, node.OutputM{
					Node:    update.Node.String(),
					Channel: string(node.ChannelActivity),
					Message: fmt.Sprintf("'%s' has reached status '%s': %s\n", n.Name(), update.Status, update.Message),
				})
			}
		}
	}()
	go func() {
		for output := range outputChan {
			runtime.EventsEmit(a.ctx, EventWorkflowOutput, output.Pack())
		}
	}()
	runtime.EventsEmit(a.ctx, EventWorkflowStarted, w.ID)
	defer runtime.EventsEmit(a.ctx, EventWorkflowFinished, w.ID)
	if err := flow.Run(ctx, updateChan, outputChan); err != nil {
		if errors.Is(err, context.Canceled) {
			a.Notify("Workflow canceled", "The workflow was canceled.")
		} else {
			a.Error("Workflow error", err.Error())
		}
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

func (a *App) CreateWorkflowFromRequest(reqU map[string]interface{}) *workflow.WorkflowM {
	j, err := json.Marshal(reqU)
	if err != nil {
		a.Error("Error creating workflow", err.Error())
		return nil
	}
	var req packaging.HttpRequest
	if err := json.Unmarshal(j, &req); err != nil {
		a.Error("Error creating workflow", err.Error())
		return nil
	}

	flow := workflow.NewWorkflow()
	reqNode := node.NewRequest()
	if err := reqNode.SetStaticInputValues(map[string]transmission.Transmission{
		"input": transmission.NewRequest(req),
	}); err != nil {
		a.Error("Error creating workflow", err.Error())
		return nil
	}
	flow.Nodes = append(flow.Nodes, reqNode)
	w, err := flow.Pack()
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
	case node.TypeVariables:
		s, err := workflow.ToNodeM(node.NewVars())
		if err != nil {
			return nil
		}
		return s
	default:
		return nil
	}
}
