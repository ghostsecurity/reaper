package workflow

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type VerifierNode struct {
	id            uuid.UUID
	input         Connector
	outputs       Connectors
	statusCodeMin int
	statusCodeMax int
}

func NewVerifierNode(statusCodeMin, statusCodeMax int) *VerifierNode {
	return &VerifierNode{
		id:    uuid.New(),
		input: NewConnector("response", TransmissionTypeRequestAndResponse),
		outputs: Connectors{
			NewConnector("good", TransmissionTypeRequestAndResponse),
			NewConnector("bad", TransmissionTypeRequestAndResponse),
		},
		statusCodeMin: statusCodeMin,
		statusCodeMax: statusCodeMax,
	}
}

func (n *VerifierNode) ID() uuid.UUID {
	return n.id
}

func (n *VerifierNode) Name() string {
	return "Verifier"
}

func (n *VerifierNode) GetInput() Connector {
	return n.input
}

func (n *VerifierNode) GetOutputs() Connectors {
	return n.outputs
}

func (n *VerifierNode) Run(ctx context.Context, in map[uuid.UUID]Transmission) (<-chan OutputInstance, <-chan error) {

	output := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(output)
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
			errs <- fmt.Errorf("input not found: no request/response specified")
			return
		}

		request := provider.Request()
		response := provider.Response()
		params := provider.ParameterSet()

		goodOutput, ok := n.outputs.FindByName("good")
		if !ok {
			errs <- fmt.Errorf("output not found: no good output specified")
			return
		}
		badOutput, ok := n.outputs.FindByName("bad")
		if !ok {
			errs <- fmt.Errorf("output not found: no bad output specified")
			return
		}

		if response.StatusCode >= n.statusCodeMin && response.StatusCode <= n.statusCodeMax {
			output <- OutputInstance{
				OutputID: goodOutput.ID,
				Current:  1,
				Total:    1,
				Complete: true,
				Data:     NewRequestAndResponseTransmission(request, response, params),
			}
		} else {
			output <- OutputInstance{
				OutputID: badOutput.ID,
				Current:  1,
				Total:    1,
				Complete: true,
				Data:     NewRequestAndResponseTransmission(request, response, params),
			}
		}
	}()

	return output, errs

}
