package node

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"golang.org/x/net/context"
)

type SenderNode struct {
	*base
	noInjections
}

func NewSender() *SenderNode {
	return &SenderNode{
		base: newBase(
			"Sender",
			TypeSender,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("start", transmission.TypeStart, true),
					NewConnector("request", transmission.TypeRequest, true),
					NewConnector("replacements", transmission.TypeMap, true),
					NewConnector("timeout", transmission.TypeInt, false, "in milliseconds"),
					NewConnector("follow_redirects", transmission.TypeBoolean, false, ""),
					NewConnector("parallelism", transmission.TypeInt, false, "number of parallel requests"),
				},
				Connectors{
					NewConnector("output", transmission.TypeRequest|transmission.TypeResponse|transmission.TypeMap, true),
				},
				map[string]transmission.Transmission{
					"timeout":          transmission.NewInt(5000),
					"follow_redirects": transmission.NewBoolean(false),
					"parallelism":      transmission.NewInt(1),
				},
			),
		),
	}
}

func (n *SenderNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {

	timeout, err := n.ReadInputInt("timeout", nil)
	if err != nil {
		return err
	}

	parallel, err := n.ReadInputInt("parallelism", nil)
	if err != nil {
		return err
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = parallel
	t.MaxConnsPerHost = parallel
	t.MaxIdleConnsPerHost = parallel

	client := http.Client{
		CheckRedirect: nil,
		Timeout:       time.Millisecond * time.Duration(timeout),
	}

	follow, err := n.ReadInputBool("follow_redirects", nil)
	if err != nil {
		return err
	}

	if !follow {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	if parallel < 1 {
		return fmt.Errorf("parallelism must be greater than 0")
	}

	restrict := make(chan struct{}, parallel)

	var wg sync.WaitGroup
	defer wg.Wait()

	defer n.setBusy(false)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case input, ok := <-in:
			if !ok {
				return nil
			}
			if input.Data == nil {
				return fmt.Errorf("input is nil")
			}

			request, err := n.ReadInputRequest("request", input.Data)
			if err != nil {
				return err
			}

			replacements, _ := n.ReadInputMap("replacements", input.Data)
			for k, v := range replacements {
				request.URL = strings.ReplaceAll(request.URL, k, v)
				request.Body = strings.ReplaceAll(request.Body, k, v)
				for i, header := range request.Headers {
					request.Headers[i].Value = strings.ReplaceAll(header.Value, k, v)
				}
			}

			r, err := packaging.UnpackageHttpRequest(request)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case restrict <- struct{}{}:
			}

			n.setBusy(true)
			wg.Add(1)

			go func() {
				defer func() {
					<-restrict
					if len(restrict) == 0 {
						n.setBusy(false)
					}
				}()
				defer wg.Done()

				resp, err := client.Do(r)
				if err != nil {
					return
				}

				response, err := packaging.PackageHttpResponse(resp, "", 0)
				if err != nil {
					return
				}

				n.tryOut(ctx, out, OutputInstance{
					OutputName: "output",
					Complete:   input.Last,
					Data:       transmission.NewRequestResponsePairWithMap(*request, *response, replacements),
				})
			}()
		}
	}
}
