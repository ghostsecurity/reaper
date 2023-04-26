package node

import (
	"fmt"
	"io"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type VerifierNode struct {
	id      uuid.UUID
	storage VarStorage
	vars    *VarStorage
}

func NewVerifier() *VerifierNode {
	return &VerifierNode{
		id: uuid.New(),
		vars: NewVarStorage(
			Connectors{
				NewConnector("response", transmission.TypeRequest|transmission.TypeResponse),
				NewConnector("min", transmission.TypeInt),
				NewConnector("max", transmission.TypeInt),
			},
			Connectors{
				NewConnector("good", transmission.TypeRequest|transmission.TypeResponse),
				NewConnector("bad", transmission.TypeRequest|transmission.TypeResponse),
			},
		),
	}
}

func (n *VerifierNode) ID() uuid.UUID {
	return n.id
}

func (n *VerifierNode) Name() string {
	return "Verifier"
}

func (n *VerifierNode) GetInputs() Connectors {
	return n.vars.GetInputs()
}

func (n *VerifierNode) GetOutputs() Connectors {
	return n.vars.GetOutputs()
}

func (n *VerifierNode) SetStaticInputValues(values map[string]transmission.Transmission) error {
	return n.vars.SetStaticInputValues(values)
}

func (n *VerifierNode) Type() NodeType {
	return TypeVerifier
}

func (n *VerifierNode) Validate(params map[string]transmission.Transmission) error {
	return n.vars.Validate(params)
}

func (n *VerifierNode) GetVars() *VarStorage {
	return n.vars
}

func (n *VerifierNode) SetVars(v *VarStorage) {
	n.vars = v
}

func (n *VerifierNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *VerifierNode) Run(ctx context.Context, in map[string]transmission.Transmission, _ io.Writer) (<-chan OutputInstance, <-chan error) {

	output := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(output)
		defer close(errs)
		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}
		response, err := n.vars.ReadInputResponse("response", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no response specified: %w", err)
			return
		}

		min, err := n.vars.ReadInputInt("min", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no min specified: %w", err)
			return
		}

		max, err := n.vars.ReadInputInt("max", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no max specified: %w", err)
			return
		}

		input, err := n.vars.ReadValue("response", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no response specified: %w", err)
			return
		}

		if response.StatusCode >= min && response.StatusCode <= max {
			output <- OutputInstance{
				OutputName: "good",
				Current:    1,
				Total:      1,
				Complete:   true,
				Data:       input,
			}
		} else {
			output <- OutputInstance{
				OutputName: "bad",
				Current:    1,
				Total:      1,
				Complete:   true,
				Data:       input,
			}
		}
	}()

	return output, errs
}
