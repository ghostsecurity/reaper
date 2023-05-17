package node

import (
	"github.com/ghostsecurity/reaper/backend/packaging"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type RequestNode struct {
	*base
}

func NewRequest() *RequestNode {
	return &RequestNode{
		base: newBase(
			"Request",
			TypeRequest,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("input", transmission.TypeRequest, false),
				},
				Connectors{
					NewConnector("output", transmission.TypeRequest, true),
				},
				map[string]transmission.Transmission{
					"input": transmission.NewRequest(packaging.HttpRequest{
						Method:      "GET",
						URL:         "https://example.com/",
						Host:        "example.com",
						Path:        "/",
						QueryString: "",
						Scheme:      "https",
						Body:        "",
						Headers: []packaging.KeyValue{
							{
								Key:   "Host",
								Value: "example.com",
							},
							{
								Key:   "User-Agent",
								Value: "Reaper",
							},
						},
						Query: []packaging.KeyValue{},
						Tags:  []string{},
					}),
				},
			),
		),
	}
}

func (n *RequestNode) GetInjections() map[string]transmission.Transmission {
	req, err := n.ReadInputRequest("input", nil)
	if err != nil {
		return nil
	}
	return map[string]transmission.Transmission{
		"output": transmission.NewRequest(*req),
	}
}

func (n *RequestNode) Run(ctx context.Context, in map[string]transmission.Transmission, output chan<- Output, last bool) (<-chan OutputInstance, <-chan error) {
	out := make(chan OutputInstance)
	errs := make(chan error)
	go func() {
		defer close(out)
		defer close(errs)
		req, err := n.ReadInputRequest("input", nil)
		if err != nil {
			errs <- err
			return
		}
		out <- OutputInstance{
			OutputName: "output",
			Data:       transmission.NewRequest(*req),
			Complete:   last,
		}
	}()
	return out, errs
}
