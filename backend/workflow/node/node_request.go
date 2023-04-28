package node

import (
	"io"

	"github.com/ghostsecurity/reaper/backend/packaging"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type RequestNode struct {
	id   uuid.UUID
	name string
	*VarStorage
}

func NewRequest() *RequestNode {
	return &RequestNode{
		id: uuid.New(),
		VarStorage: NewVarStorage(
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
					Tags:  nil,
				}),
			},
		),
		name: "Request",
	}
}

func (n *RequestNode) ID() uuid.UUID {
	return n.id
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

func (n *RequestNode) Name() string {
	return n.name
}

func (n *RequestNode) SetName(name string) {
	n.name = name
}

func (n *RequestNode) Type() Type {
	return TypeRequest
}

func (n *RequestNode) IsReadOnly() bool {
	return false
}

func (n *RequestNode) GetVars() *VarStorage {
	return n.VarStorage
}

func (n *RequestNode) SetVars(v *VarStorage) {
	n.VarStorage = v
}

func (n *RequestNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *RequestNode) Run(ctx context.Context, in map[string]transmission.Transmission, _, _ io.Writer) (<-chan OutputInstance, <-chan error) {
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
			Complete:   true,
		}
	}()
	return out, errs
}
