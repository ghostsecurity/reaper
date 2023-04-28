package node

import (
	"fmt"
	"io"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type FuzzerNode struct {
	*VarStorage
	id   uuid.UUID
	name string
}

func NewFuzzer() *FuzzerNode {
	return &FuzzerNode{
		id:   uuid.New(),
		name: "Fuzzer",
		VarStorage: NewVarStorage(
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
	}
}

func (n *FuzzerNode) IsReadOnly() bool {
	return false
}

func (n *FuzzerNode) ID() uuid.UUID {
	return n.id
}

func (n *FuzzerNode) Name() string {
	return n.name
}

func (n *FuzzerNode) Type() Type {
	return TypeFuzzer
}

func (n *FuzzerNode) SetName(name string) {
	n.name = name
}

func (n *FuzzerNode) GetInjections() map[string]transmission.Transmission {
	return nil
}

func (n *FuzzerNode) GetVars() *VarStorage {
	return n.VarStorage
}

func (n *FuzzerNode) SetVars(vars *VarStorage) {
	n.VarStorage = vars
}

func (n *FuzzerNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *FuzzerNode) Run(ctx context.Context, in map[string]transmission.Transmission, _, _ io.Writer) (<-chan OutputInstance, <-chan error) {

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

			data := map[string]string{
				placeholder: word,
			}
			for k, v := range vars {
				data[k] = v
			}

			output <- OutputInstance{
				OutputName: "output",
				Current:    int(i),
				Total:      list.Count(),
				Complete:   list.Complete(),
				Data:       transmission.NewMap(data),
			}
		}
	}()

	return output, errs

}
