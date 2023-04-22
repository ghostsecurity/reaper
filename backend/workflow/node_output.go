package workflow

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type OutputNode struct {
	id    uuid.UUID
	name  string
	w     io.Writer
	input Connector
}

func NewOutputNode(w io.Writer, name string) *OutputNode {
	return &OutputNode{
		id:    uuid.New(),
		w:     w,
		input: NewConnector("input", TransmissionTypeRequestAndResponse),
		name:  name,
	}
}

func (n *OutputNode) ID() uuid.UUID {
	return n.id
}

func (n *OutputNode) Name() string {
	return n.name
}

func (n *OutputNode) GetInput() Connector {
	return n.input
}

func (n *OutputNode) GetOutputs() Connectors {
	return nil
}

func (n *OutputNode) Run(ctx context.Context, in map[uuid.UUID]Transmission) (<-chan OutputInstance, <-chan error) {

	out := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(out)
		defer close(errs)
		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}
		transmission, ok := in[n.input.ID]
		if !ok {
			errs <- fmt.Errorf("input not found: no request/response pair specified")
			return
		}
		provider, ok := transmission.(RequestAndResponseProvider)
		if !ok {
			errs <- fmt.Errorf("input was not expected type: '%#v'", transmission)
			return
		}
		req, resp, params := provider.Request(), provider.Response(), provider.ParameterSet()

		var paramSummary []string
		for k, v := range params {
			paramSummary = append(paramSummary, fmt.Sprintf("%s=%s", k, v))
		}

		fmt.Fprintf(n.w, "%s %s -> %d (%s)\n", req.Method, req.URL, resp.StatusCode, strings.Join(paramSummary, ", "))
	}()

	return out, errs
}
