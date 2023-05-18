package node

import (
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type VarsNode struct {
	*base
}

func NewVars() *VarsNode {
	return &VarsNode{
		base: newBase(
			"Variables",
			TypeVariables,
			false,
			NewVarStorage(
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
		),
	}
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

func (n *VarsNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {
	return nil
}
