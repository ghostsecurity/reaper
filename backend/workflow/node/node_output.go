package node

import (
	"fmt"
	"io"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type OutputNode struct {
	id   uuid.UUID
	name string
	*VarStorage
}

func NewOutput() *OutputNode {
	return &OutputNode{
		id: uuid.New(),
		VarStorage: NewVarStorage(
			Connectors{
				NewConnector("input", transmission.TypeAny, true),
				NewConnector("stdout", transmission.TypeBoolean, false),
				NewConnector("stderr", transmission.TypeBoolean, false),
			},
			nil,
			map[string]transmission.Transmission{
				"stdout": transmission.NewBoolean(true),
				"stderr": transmission.NewBoolean(false),
			},
		),
		name: "Output",
	}
}

func (n *OutputNode) ID() uuid.UUID {
	return n.id
}

func (n *OutputNode) GetInjections() map[string]transmission.Transmission {
	return nil
}

func (n *OutputNode) Name() string {
	return n.name
}

func (n *OutputNode) SetName(name string) {
	n.name = name
}

func (n *OutputNode) Type() Type {
	return TypeOutput
}

func (n *OutputNode) GetVars() *VarStorage {
	return n.VarStorage
}

func (n *OutputNode) SetVars(v *VarStorage) {
	n.VarStorage = v
}

func (n *OutputNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *OutputNode) IsReadOnly() bool {
	return false
}

func (n *OutputNode) Run(ctx context.Context, in map[string]transmission.Transmission, stdout, stderr io.Writer) (<-chan OutputInstance, <-chan error) {

	isOut, _ := n.ReadInputBool("stdout", nil)
	isErr, _ := n.ReadInputBool("stderr", nil)

	if !isOut && !isErr {
		return nil, nil
	}

	printf := func(format string, args ...interface{}) {
		if isOut {
			fmt.Fprintf(stdout, format, args...)
		}
		if isErr {
			fmt.Fprintf(stderr, format, args...)
		}
	}

	out := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(out)
		defer close(errs)

		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}
		if req, err := n.ReadInputRequest("input", in); err == nil {
			printf("%s %s", req.Method, req.URL)
			if resp, err := n.ReadInputResponse("input", in); err == nil {
				printf(" -> %d", resp.StatusCode)
			}
		}

		if params, err := n.ReadInputMap("input", in); err == nil {
			if len(params) > 0 {
				printf(" [")
				var i int
				for k, v := range params {
					printf("%s=%s", k, v)
					if i < len(params)-1 {
						printf(" ")
					}
				}
				printf("]")
			}
		}

		printf("\n")
	}()

	return out, errs
}
