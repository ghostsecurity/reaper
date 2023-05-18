package node

import (
	"fmt"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type StatusFilterNode struct {
	*base
	noInjections
}

func NewStatusFilter() *StatusFilterNode {
	return &StatusFilterNode{
		base: newBase(
			"Status Filter",
			TypeStatusFilter,
			false,
			NewVarStorage(
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
		),
	}
}

func (n *StatusFilterNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {

	defer n.setBusy(false)

	min, err := n.ReadInputInt("min", nil)
	if err != nil {
		return fmt.Errorf("input not found: no min specified: %w", err)
	}

	max, err := n.ReadInputInt("max", nil)
	if err != nil {
		return fmt.Errorf("input not found: no max specified: %w", err)
	}

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
				return fmt.Errorf("input data is nil")
			}

			response, err := n.ReadInputResponse("response", input.Data)
			if err != nil {
				return fmt.Errorf("input not found: no response specified: %w", err)
			}

			rawInput, err := n.ReadValue("response", input.Data)
			if err != nil {
				return fmt.Errorf("input not found: no response specified: %w", err)
			}

			if response.StatusCode >= min && response.StatusCode <= max {
				n.tryOut(ctx, out, OutputInstance{
					OutputName: "good",
					Current:    1,
					Total:      1,
					Complete:   input.Last,
					Data:       rawInput,
				})
			} else {
				n.tryOut(ctx, out, OutputInstance{
					OutputName: "bad",
					Current:    1,
					Total:      1,
					Complete:   input.Last,
					Data:       rawInput,
				})
			}

			n.setBusy(false)

		}
	}

}
