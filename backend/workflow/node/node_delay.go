package node

import (
	"fmt"
	"time"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type DelayNode struct {
	*base
	noInjections
}

func NewDelay() *DelayNode {
	return &DelayNode{
		base: newBase(
			"Delay",
			TypeDelay,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("input", transmission.TypeAny, true),
					NewConnector("delay", transmission.TypeInt, false, "in milliseconds"),
				},
				Connectors{
					NewConnector("output", transmission.TypeAny, true),
				},
				map[string]transmission.Transmission{
					"delay": transmission.NewInt(1000),
				},
			),
		),
	}
}

func (n *DelayNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {

	delay, err := n.ReadInputInt("delay", nil)
	if err != nil {
		return err
	}

	defer n.setBusy(false)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case input, ok := <-in:
			if !ok {
				return nil
			}

			n.setBusy(true)

			if input.Data == nil {
				return fmt.Errorf("input is nil")
			}

			raw, err := n.ReadValue("input", input.Data)
			if err != nil {
				return fmt.Errorf("input not found: %v", err)
			}

			time.Sleep(time.Duration(delay) * time.Millisecond)

			n.tryOut(ctx, out, OutputInstance{
				OutputName: "output",
				Complete:   input.Last,
				Data:       raw,
			})

			n.setBusy(false)
		}
	}
}
