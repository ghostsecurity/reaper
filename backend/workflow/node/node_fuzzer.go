package node

import (
	"fmt"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type FuzzerNode struct {
	*base
	noInjections
}

func NewFuzzer() *FuzzerNode {
	return &FuzzerNode{
		base: newBase(
			"Fuzzer",
			TypeFuzzer,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("start", transmission.TypeStart, true),
					NewConnector("vars", transmission.TypeMap, true),
					NewConnector("placeholder", transmission.TypeString, false),
					NewConnector("list", transmission.TypeList, false),
				},
				Connectors{
					NewConnector("output", transmission.TypeMap, true),
				},
				map[string]transmission.Transmission{
					"placeholder": transmission.NewString("$FUZZ$"),
					"list":        transmission.NewNumericRangeIterator(0, 100),
				},
			),
		),
	}
}

func (n *FuzzerNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {

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
			list, err := n.ReadInputList("list", input.Data)
			if err != nil {
				return fmt.Errorf("input not found: no list specified")
			}
			placeholder, err := n.ReadInputString("placeholder", input.Data)
			if err != nil {
				return fmt.Errorf("input not found: no placeholder specified")
			}

			vars, _ := n.ReadInputMap("vars", input.Data)

			var i int64
			for {
				i++
				word, ok := list.Next()
				if !ok {
					break
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				data := map[string]string{}
				for k, v := range vars {
					data[k] = v
				}
				data[placeholder] = word

				n.tryOut(ctx, out, OutputInstance{
					OutputName: "output",
					Current:    int(i),
					Total:      list.Count(),
					Complete:   list.Complete() && input.Last,
					Data:       transmission.NewMap(data),
				})

			}

			n.setBusy(false)
		}
	}
}
