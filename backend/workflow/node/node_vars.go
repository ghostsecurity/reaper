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
