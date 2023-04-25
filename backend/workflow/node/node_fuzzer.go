package node

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type FuzzerNode struct {
	id     uuid.UUID
	vars   *VarStorage
	client *http.Client
}

func NewFuzzer() *FuzzerNode {
	return &FuzzerNode{
		id: uuid.New(),
		vars: NewVarStorage(
			Connectors{
				NewConnector("request", transmission.TypeRequest),
				NewConnector("placeholder", transmission.TypeString),
				NewConnector("list", transmission.TypeList),
			},
			Connectors{
				NewConnector("output", transmission.TypeRequest|transmission.TypeResponse),
			},
		),
		client: http.DefaultClient,
	}
}

func (n *FuzzerNode) ID() uuid.UUID {
	return n.id
}

func (n *FuzzerNode) Name() string {
	return "Fuzzer"
}

func (n *FuzzerNode) Type() NodeType {
	return TypeFuzzer
}

func (n *FuzzerNode) GetInputs() Connectors {
	return n.vars.GetInputs()
}

func (n *FuzzerNode) GetOutputs() Connectors {
	return n.vars.GetOutputs()
}

func (n *FuzzerNode) SetStaticInputValues(values map[string]transmission.Transmission) {
	n.vars.SetStaticInputValues(values)
}

func (n *FuzzerNode) GetVars() *VarStorage {
	return n.vars
}

func (n *FuzzerNode) SetVars(vars *VarStorage) {
	n.vars = vars
}

func (n *FuzzerNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *FuzzerNode) Run(ctx context.Context, in map[string]transmission.Transmission) (<-chan OutputInstance, <-chan error) {

	output := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(output)
		defer close(errs)
		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}
		list, err := n.vars.ReadInputList("list", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no list specified")
			return
		}
		placeholder, err := n.vars.ReadInputString("placeholder", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no placeholder specified")
			return
		}

		var i int64
		for {
			i++
			word, ok := list.Next()
			if !ok {
				break
			}
			request, err := n.vars.ReadInputRequest("request", in)
			if err != nil {
				errs <- fmt.Errorf("input not found: no request specified")
				return
			}
			request.URL = strings.ReplaceAll(request.URL, placeholder, word)
			request.Body = strings.ReplaceAll(request.Body, placeholder, word)
			for i, header := range request.Headers {
				request.Headers[i].Value = strings.ReplaceAll(header.Value, placeholder, word)
			}

			r, err := packaging.UnpackageHttpRequest(request)
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
				OutputName: "output",
				Current:    int(i),
				Total:      list.Count(),
				Complete:   list.Complete(),
				Data: transmission.NewRequestResponsePairWithMap(*request, *response, map[string]string{
					placeholder: word,
				}),
			}
		}
	}()

	return output, errs

}
