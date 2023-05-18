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

func (n *StartNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {
	defer n.setBusy(false)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case input, ok := <-in:
		if !ok {
			return nil
		}

		n.tryOut(ctx, out, OutputInstance{
			OutputName: "output",
			Data:       transmission.NewStart(),
			Complete:   input.Last,
		})

	}
	return nil
}
