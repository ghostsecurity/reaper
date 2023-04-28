package node

import (
	"fmt"
	"io"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type StatusFilterNode struct {
	id   uuid.UUID
	name string
	*VarStorage
}

func NewStatusFilter() *StatusFilterNode {
	return &StatusFilterNode{
		id:   uuid.New(),
		name: "Status Filter",
		VarStorage: NewVarStorage(
			Connectors{
				NewConnector("response", transmission.TypeRequest|transmission.TypeResponse|transmission.TypeMap, true),
				NewConnector("min", transmission.TypeInt, false),
				NewConnector("max", transmission.TypeInt, false),
			},
			Connectors{
				NewConnector("good", transmission.TypeRequest|transmission.TypeResponse|transmission.TypeMap, true),
				NewConnector("bad", transmission.TypeRequest|transmission.TypeResponse|transmission.TypeMap, true),
			},
			map[string]transmission.Transmission{
				"min": transmission.NewInt(200),
				"max": transmission.NewInt(399),
			},
		),
	}
}

func (n *StatusFilterNode) ID() uuid.UUID {
	return n.id
}

func (n *StatusFilterNode) IsReadOnly() bool {
	return false
}

func (n *StatusFilterNode) Name() string {
	return n.name
}

func (n *StatusFilterNode) SetName(name string) {
	n.name = name
}

func (n *StatusFilterNode) GetInjections() map[string]transmission.Transmission {
	return nil
}

func (n *StatusFilterNode) Type() Type {
	return TypeStatusFilter
}

func (n *StatusFilterNode) GetVars() *VarStorage {
	return n.VarStorage
}

func (n *StatusFilterNode) SetVars(v *VarStorage) {
	n.VarStorage = v
}

func (n *StatusFilterNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *StatusFilterNode) Run(ctx context.Context, in map[string]transmission.Transmission, _, _ io.Writer) (<-chan OutputInstance, <-chan error) {

	output := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(output)
		defer close(errs)
		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}
		response, err := n.ReadInputResponse("response", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no response specified: %w", err)
			return
		}

		min, err := n.ReadInputInt("min", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no min specified: %w", err)
			return
		}

		max, err := n.ReadInputInt("max", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no max specified: %w", err)
			return
		}

		input, err := n.ReadValue("response", in)
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
