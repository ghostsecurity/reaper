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

func (n *FuzzerNode) Run(ctx context.Context, in map[string]transmission.Transmission, out chan<- Output, last bool) (<-chan OutputInstance, <-chan error) {

	output := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(output)
		defer close(errs)
		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}
		list, err := n.ReadInputList("list", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no list specified")
			return
		}
		placeholder, err := n.ReadInputString("placeholder", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no placeholder specified")
			return
		}

		vars, _ := n.ReadInputMap("vars", in)

		var i int64
		for {
			i++
			word, ok := list.Next()
			if !ok {
				break
			}
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
			}

			data := map[string]string{}
			for k, v := range vars {
				data[k] = v
			}
			data[placeholder] = word

			output <- OutputInstance{
				OutputName: "output",
				Current:    int(i),
				Total:      list.Count(),
				Complete:   list.Complete() && last,
				Data:       transmission.NewMap(data),
			}
		}
	}()

	return output, errs

}
