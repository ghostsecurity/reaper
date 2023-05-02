package node

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type SenderNode struct {
	*VarStorage
	id   uuid.UUID
	name string
}

func NewSender() *SenderNode {
	return &SenderNode{
		id:   uuid.New(),
		name: "Sender",
		VarStorage: NewVarStorage(
			Connectors{
				NewConnector("start", transmission.TypeStart, true),
				NewConnector("request", transmission.TypeRequest, true),
				NewConnector("replacements", transmission.TypeMap, true),
				NewConnector("timeout", transmission.TypeInt, false, "in milliseconds"),
				NewConnector("follow_redirects", transmission.TypeBoolean, false, ""),
			},
			Connectors{
				NewConnector("output", transmission.TypeRequest|transmission.TypeResponse|transmission.TypeMap, true),
			},
			map[string]transmission.Transmission{
				"timeout":          transmission.NewInt(5000),
				"follow_redirects": transmission.NewBoolean(false),
			},
		),
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

func (n *SenderNode) Run(ctx context.Context, in map[string]transmission.Transmission, out chan<- Output, last bool) (<-chan OutputInstance, <-chan error) {

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

		timeout, err := n.ReadInputInt("timeout", in)
		if err != nil {
			errs <- err
			return
		}

		client := http.Client{
			CheckRedirect: nil,
			Timeout:       time.Millisecond * time.Duration(timeout),
		}

		follow, err := n.ReadInputBool("follow_redirects", in)
		if err != nil {
			errs <- err
			return
		}

		if !follow {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}

		resp, err := client.Do(r)
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
			Complete:   last,
			Data:       transmission.NewRequestResponsePairWithMap(*request, *response, replacements),
		}
	}()

	return output, errs

}
