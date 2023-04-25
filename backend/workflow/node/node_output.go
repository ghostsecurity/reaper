package node

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type OutputNode struct {
	id   uuid.UUID
	name string
	w    *bytes.Buffer
	vars *VarStorage
}

func NewOutput() *OutputNode {
	return &OutputNode{
		id: uuid.New(),
		w:  bytes.NewBuffer(nil),
		vars: NewVarStorage(Connectors{
			NewConnector("input", transmission.TypeRequest|transmission.TypeResponse),
		}, nil),
		name: "output",
	}
}

func (n *OutputNode) ID() uuid.UUID {
	return n.id
}

func (n *OutputNode) Name() string {
	return n.name
}

func (n *OutputNode) GetInputs() Connectors {
	return n.vars.GetInputs()
}

func (n *OutputNode) GetOutputs() Connectors {
	return n.vars.GetOutputs()
}

func (n *OutputNode) SetStaticInputValues(values map[string]transmission.Transmission) {
	n.vars.SetStaticInputValues(values)
}

func (n *OutputNode) Type() NodeType {
	return TypeOutput
}

func (n *OutputNode) GetVars() *VarStorage {
	return n.vars
}

func (n *OutputNode) SetVars(v *VarStorage) {
	n.vars = v
}

func (n *OutputNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *OutputNode) Run(ctx context.Context, in map[string]transmission.Transmission) (<-chan OutputInstance, <-chan error) {

	out := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(out)
		defer close(errs)

		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}
		req, err := n.vars.ReadInputRequest("input", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no input: %w", err)
			return
		}

		var params map[string]string
		if m, err := n.vars.ReadInputMap("input", in); err == nil {
			params = m
		}

		var statusCode int
		if i, err := n.vars.ReadInputInt("input", in); err == nil {
			statusCode = int(i)
		}

		var paramSummary []string
		for k, v := range params {
			paramSummary = append(paramSummary, fmt.Sprintf("%s=%s", k, v))
		}

		fmt.Fprintf(n.w, "%s %s -> %d (%s)\n", req.Method, req.URL, statusCode, strings.Join(paramSummary, ", "))
	}()

	return out, errs
}

func (n *OutputNode) String() string {
	return n.w.String()
}
