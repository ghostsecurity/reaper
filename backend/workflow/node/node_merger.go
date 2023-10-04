package node

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
)

type MergerNode struct {
	*base
	noInjections
}

func NewMerger() *MergerNode {
	return &MergerNode{
		base: newBase(
			"Merger",
			TypeMerger,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("start", transmission.TypeStart, true),
					NewConnector("vars_1", transmission.TypeMap, true),
					NewConnector("vars_2", transmission.TypeMap, true),
				},
				Connectors{
					NewConnector("output", transmission.TypeMap, true),
				},
				nil,
			),
		),
	}
}

func (n *MergerNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {

	var v1s []map[string]string
	var v2s []map[string]string

	defer n.setBusy(false)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case input, ok := <-in:
			if !ok {
				return nil
			}
			if input.Data == nil {
				return fmt.Errorf("input is nil")
			}

			if vm1, err := n.ReadInputMap("vars_1", input.Data); err == nil {
				v1s = append(v1s, vm1)
			}
			if vm2, err := n.ReadInputMap("vars_2", input.Data); err == nil {
				v2s = append(v2s, vm2)
			}

			if len(v1s) == 0 || len(v2s) == 0 {
				continue
			}

			v1 := v1s[0]
			v2 := v2s[0]

			for k, v := range v2 {
				v1[k] = v
			}
			v1s = v1s[1:]
			v2s = v2s[1:]

			n.tryOut(ctx, out, OutputInstance{
				OutputName: "output",
				Complete:   input.Last,
				Data:       transmission.NewMap(v1),
			})
		}
	}
}
