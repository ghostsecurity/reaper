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

func (n *StatusFilterNode) Run(ctx context.Context, in map[string]transmission.Transmission, out chan<- Output, last bool) (<-chan OutputInstance, <-chan error) {

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
				Complete:   last,
				Data:       input,
			}
		} else {
			output <- OutputInstance{
				OutputName: "bad",
				Current:    1,
				Total:      1,
				Complete:   last,
				Data:       input,
			}
		}
	}()

	return output, errs
}
