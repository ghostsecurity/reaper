package node

import (
	"fmt"
	"strings"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type OutputNode struct {
	*base
	noInjections
}

func NewOutput() *OutputNode {
	return &OutputNode{
		base: newBase(
			"Output",
			TypeOutput,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("input", transmission.TypeAny, true),
					NewConnector("stdout", transmission.TypeBoolean, false),
					NewConnector("stderr", transmission.TypeBoolean, false),
					NewConnector("template", transmission.TypeString, false),
				},
				nil,
				map[string]transmission.Transmission{
					"stdout":   transmission.NewBoolean(true),
					"stderr":   transmission.NewBoolean(false),
					"template": transmission.NewString(""),
				},
			),
		),
	}
}

func (n *OutputNode) Start(ctx context.Context, in <-chan Input, _ chan<- OutputInstance, output chan<- Output) error {

	isOut, _ := n.ReadInputBool("stdout", nil)
	isErr, _ := n.ReadInputBool("stderr", nil)

	printf := func(format string, args ...interface{}) {
		if isOut {
			output <- Output{
				Node:    n.ID(),
				Channel: ChannelStdout,
				Message: fmt.Sprintf(format, args...),
			}
		}
		if isErr {
			output <- Output{
				Node:    n.ID(),
				Channel: ChannelStderr,
				Message: fmt.Sprintf(format, args...),
			}
		}
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

			if !isOut && !isErr {
				continue
			}

			if template, err := n.ReadInputString("template", input.Data); err == nil && template != "" {
				if params, err := n.ReadInputMap("input", input.Data); err == nil {
					for key, val := range params {
						template = strings.ReplaceAll(template, key, val)
					}
				}
				printf("%s\n", template)
				continue
			}

			if req, err := n.ReadInputRequest("input", input.Data); err == nil {
				printf("%s %s", req.Method, req.URL)
				if resp, err := n.ReadInputResponse("input", input.Data); err == nil {
					printf(" -> %d", resp.StatusCode)
				}
			}

			if params, err := n.ReadInputMap("input", input.Data); err == nil {
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

			n.setBusy(false)
		}
	}

}
