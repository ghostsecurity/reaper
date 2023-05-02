package node

import (
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type VarsNode struct {
	id   uuid.UUID
	name string
	*VarStorage
}

func NewVars() *VarsNode {
	return &VarsNode{
		id: uuid.New(),
		VarStorage: NewVarStorage(
			Connectors{
				NewConnector("variables", transmission.TypeMap, false),
			},
			Connectors{
				NewConnector("output", transmission.TypeMap, true),
			},
			map[string]transmission.Transmission{
				"variables": transmission.NewMap(map[string]string{}),
			},
		),
		name: "Variables",
	}
}

func (n *VarsNode) ID() uuid.UUID {
	return n.id
}

func (n *VarsNode) GetInjections() map[string]transmission.Transmission {
	m, err := n.ReadInputMap("variables", nil)
	if err != nil {
		return nil
	}
	return map[string]transmission.Transmission{
		"output": transmission.NewMap(m),
	}
}

func (n *VarsNode) Name() string {
	return n.name
}

func (n *VarsNode) SetName(name string) {
	n.name = name
}

func (n *VarsNode) Type() Type {
	return TypeVariables
}

func (n *VarsNode) IsReadOnly() bool {
	return false
}

func (n *VarsNode) GetVars() *VarStorage {
	return n.VarStorage
}

func (n *VarsNode) SetVars(v *VarStorage) {
	n.VarStorage = v
}

func (n *VarsNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *VarsNode) Run(ctx context.Context, in map[string]transmission.Transmission, _ chan<- Output, last bool) (<-chan OutputInstance, <-chan error) {
	out := make(chan OutputInstance)
	errs := make(chan error)
	go func() {
		defer close(out)
		defer close(errs)
		m, err := n.ReadInputMap("variables", nil)
		if err != nil {
			errs <- err
			return
		}
		out <- OutputInstance{
			OutputName: "output",
			Data:       transmission.NewMap(m),
			Complete:   last,
		}
	}()
	return out, errs
}
