package workflow

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type FuzzerNode struct {
	id          uuid.UUID
	input       Connector
	output      Connector
	list        StringInterator
	placeholder string
	client      *http.Client
}

type StringInterator interface {
	Next() (string, bool)
	Count() int
	Complete() bool
}

func NewFuzzerNode(placeholder string, list StringInterator) *FuzzerNode {
	return &FuzzerNode{
		id:          uuid.New(),
		placeholder: placeholder,
		list:        list,
		input:       NewConnector("request", TransmissionTypeRequest),
		output:      NewConnector("request", TransmissionTypeRequestAndResponse),
		client:      http.DefaultClient,
	}
}

func (n *FuzzerNode) ID() uuid.UUID {
	return n.id
}

func (n *FuzzerNode) Name() string {
	return "Fuzzer"
}

func (n *FuzzerNode) GetInput() Connector {
	return n.input
}

func (n *FuzzerNode) GetOutputs() Connectors {
	return Connectors{n.output}
}

func (n *FuzzerNode) Run(ctx context.Context, in map[uuid.UUID]Transmission) (<-chan OutputInstance, <-chan error) {

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
			errs <- fmt.Errorf("input not found: no request specified")
			return
		}
		requestProvider, ok := transmission.(RequestProvider)
		if !ok {
			errs <- fmt.Errorf("input not found: no request specified")
			return
		}
		var i int64
		for {
			i++
			word, ok := n.list.Next()
			if !ok {
				break
			}
			request := requestProvider.Request()
			request.URL = strings.ReplaceAll(request.URL, n.placeholder, word)
			request.Body = strings.ReplaceAll(request.Body, n.placeholder, word)
			for i, header := range request.Headers {
				request.Headers[i].Value = strings.ReplaceAll(header.Value, n.placeholder, word)
			}

			r, err := packaging.UnpackageHttpRequest(&request)
			if err != nil {
				errs <- err
				return
			}

			resp, err := n.client.Do(r)
			if err != nil {
				errs <- err
				return
			}

			response, err := packaging.PackageHttpResponse(resp, n.id.String(), i)
			if err != nil {
				errs <- err
				return
			}

			output <- OutputInstance{
				OutputID: n.output.ID,
				Current:  int(i),
				Total:    n.list.Count(),
				Complete: n.list.Complete(),
				Data: NewRequestAndResponseTransmission(request, *response, map[string]string{
					n.placeholder: word,
				}),
			}
		}
	}()

	return output, errs

}
