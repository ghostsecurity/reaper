package node

import (
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type StartNode struct {
	*base
	noInjections
}

func NewStart() *StartNode {
	return &StartNode{
		base: newBase(
			"Start",
			TypeStart,
			true,
			NewVarStorage(
				nil,
				Connectors{
					NewConnector("output", transmission.TypeStart, true),
				},
				nil,
			),
		),
	}
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
