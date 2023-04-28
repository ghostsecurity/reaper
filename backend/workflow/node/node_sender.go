package node

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type SenderNode struct {
	*VarStorage
	id     uuid.UUID
	name   string
	client *http.Client
}

func NewSender() *SenderNode {
	return &SenderNode{
		id:   uuid.New(),
		name: "Sender",
		VarStorage: NewVarStorage(
			Connectors{
				NewConnector("request", transmission.TypeRequest, true),
				NewConnector("replacements", transmission.TypeMap, true),
			},
			Connectors{
				NewConnector("output", transmission.TypeRequest|transmission.TypeResponse|transmission.TypeMap, true),
			},
			nil,
		),
		client: http.DefaultClient,
	}
}

func (n *SenderNode) IsReadOnly() bool {
	return false
}

func (n *SenderNode) ID() uuid.UUID {
	return n.id
}

func (n *SenderNode) Name() string {
	return n.name
}

func (n *SenderNode) Type() Type {
	return TypeSender
}

func (n *SenderNode) SetName(name string) {
	n.name = name
}

func (n *SenderNode) GetInjections() map[string]transmission.Transmission {
	return nil
}

func (n *SenderNode) GetVars() *VarStorage {
	return n.VarStorage
}

func (n *SenderNode) SetVars(vars *VarStorage) {
	n.VarStorage = vars
}

func (n *SenderNode) SetID(id uuid.UUID) {
	n.id = id
}

func (n *SenderNode) Run(ctx context.Context, in map[string]transmission.Transmission, _, _ io.Writer) (<-chan OutputInstance, <-chan error) {

	output := make(chan OutputInstance)
	errs := make(chan error)

	go func() {
		defer close(output)
		defer close(errs)
		if in == nil {
			errs <- fmt.Errorf("input is nil")
			return
		}

		request, err := n.ReadInputRequest("request", in)
		if err != nil {
			errs <- fmt.Errorf("input not found: no request specified")
			return
		}

		replacements, _ := n.ReadInputMap("replacements", in)
		for k, v := range replacements {
			request.URL = strings.ReplaceAll(request.URL, k, v)
			request.Body = strings.ReplaceAll(request.Body, k, v)
			for i, header := range request.Headers {
				request.Headers[i].Value = strings.ReplaceAll(header.Value, k, v)
			}
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

		response, err := packaging.PackageHttpResponse(resp, "", 0)
		if err != nil {
			errs <- err
			return
		}

		output <- OutputInstance{
			OutputName: "output",
			Complete:   false,
			Data:       transmission.NewRequestResponsePairWithMap(*request, *response, replacements),
		}
	}()

	return output, errs

}
