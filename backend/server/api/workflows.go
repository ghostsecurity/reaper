package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"

	"github.com/ghostsecurity/reaper/backend/packaging"

	"github.com/google/uuid"

	"github.com/ghostsecurity/reaper/backend/workflow"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
)

func (a *API) RunWorkflow(w *workflow.WorkflowM) {
	go a.runWorkflow(w)
}

func (a *API) runWorkflow(w *workflow.WorkflowM) {
	a.workflowMu.Lock()
	defer a.workflowMu.Unlock()
	if a.runningWorkflowID != uuid.Nil {
		a.notify("There is already a workflow running. Please cancel it or wait for it to finish.")
		return
	}
	flow, err := w.Unpack()
	if err != nil {
		a.notify("Workflow unpack error: %s", err.Error())
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
			_ = a.eventTrigger(EventWorkflowUpdate, update.Pack())
			if n, err := flow.FindNode(update.Node); err == nil {
				_ = a.eventTrigger(EventWorkflowOutput, node.OutputM{
					Node:    update.Node.String(),
					Channel: string(node.ChannelActivity),
					Message: fmt.Sprintf("'%s' has reached status '%s': %s\n", n.Name(), update.Status, update.Message),
				})
			}
		}
	}()
	go func() {
		for output := range outputChan {
			_ = a.eventTrigger(EventWorkflowOutput, output.Pack())
		}
	}()
	_ = a.eventTrigger(EventWorkflowStarted, w.ID)
	defer func() {
		_ = a.eventTrigger(EventWorkflowFinished, w.ID)
	}()
	if err := flow.Run(ctx, updateChan, outputChan); err != nil {
		if errors.Is(err, context.Canceled) {
			a.notify("Workflow canceled: %s", "The workflow was canceled.")
		} else {
			a.notify("Workflow error: %s", err.Error())
		}
	}
}

func (a *API) StopWorkflow() {
	if a.workflowContextCancel != nil {
		a.logger.Info("Stopping workflow!")
		a.workflowContextCancel()
	} else {
		a.logger.Info("No workflow to stop!")
	}
}

func (a *API) CreateWorkflow() *workflow.WorkflowM {
	w, err := workflow.NewWorkflow().Pack()
	if err != nil {
		return nil
	}
	return w
}

func (a *API) CreateWorkflowFromRequest(reqU map[string]interface{}) *workflow.WorkflowM {
	j, err := json.Marshal(reqU)
	if err != nil {
		a.notify("Error creating workflow: %s", err.Error())
		return nil
	}
	var req packaging.HttpRequest
	if err := json.Unmarshal(j, &req); err != nil {
		a.notify("Error creating workflow: %s", err.Error())
		return nil
	}

	flow := workflow.NewWorkflow()
	reqNode := node.NewRequest()
	if err := reqNode.SetStaticInputValues(map[string]transmission.Transmission{
		"input": transmission.NewRequest(req),
	}); err != nil {
		a.notify("Error creating workflow: %s", err.Error())
		return nil
	}
	flow.Nodes = append(flow.Nodes, reqNode)
	w, err := flow.Pack()
	if err != nil {
		return nil
	}
	return w
}

func (a *API) CreateNode(nodeType int) *workflow.NodeM {
	created, err := node.FromType(node.Type(nodeType))
	if err != nil {
		return nil
	}
	n, err := workflow.ToNodeM(created)
	if err != nil {
		return nil
	}
	return n
}
