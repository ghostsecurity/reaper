package node

import (
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type StartNode struct {
	id   uuid.UUID
	name string
	*VarStorage
}

func NewStart() *StartNode {
	return &StartNode{
		id: uuid.New(),
		VarStorage: NewVarStorage(
			nil,
			Connectors{
				NewConnector("output", transmission.TypeStart, true),
			},
			nil,
		),
		name: "Start",
	}
}

func (n *StartNode) ID() uuid.UUID {
	return n.id
}

func (n *StartNode) GetInjections() map[string]transmission.Transmission {
	return nil
}

func (n *StartNode) Name() string {
	return n.name
}

func (n *StartNode) SetName(name string) {
	n.name = name
}

func (n *StartNode) Type() Type {
	return TypeStart
}

func (n *StartNode) IsReadOnly() bool {
	return true
}

func (n *StartNode) GetVars() *VarStorage {
	return n.VarStorage
}

func (n *StartNode) SetVars(v *VarStorage) {
	n.VarStorage = v
}

func (n *StartNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *StartNode) Run(ctx context.Context, in map[string]transmission.Transmission, output chan<- Output, last bool) (<-chan OutputInstance, <-chan error) {
	out := make(chan OutputInstance)
	errs := make(chan error)
	go func() {
		defer close(out)
		defer close(errs)
		out <- OutputInstance{
			OutputName: "output",
			Data:       transmission.NewStart(),
			Complete:   last,
		}
	}()
	return out, errs
}
