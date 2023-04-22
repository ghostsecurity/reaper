package workflow

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FuzzingWorkflow(t *testing.T) {

	// generate random number
	secret := rand.Intn(100)

	// start server with endpoint /account?id={id} where id will only 200 with random number
	addr := ":8888"
	ln, err := net.Listen("tcp", addr)
	require.NoError(t, err)
	srv := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("id") == fmt.Sprintf("%d", secret) {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}),
	}

	go func() { _ = srv.Serve(ln) }()

	fuzzer := NewFuzzerNode("$X$", NewNumericRangeIterator(0, 100))
	verifier := NewVerifierNode(200, 200)

	output := bytes.NewBuffer(nil)

	nOutput := NewOutputNode(output, "output")
	nError := NewOutputNode(os.Stderr, "error")

	verifierGood, ok := verifier.GetOutputs().FindByName("good")
	require.True(t, ok)

	// start workflow with input node and output node
	workflow := &Workflow{
		ID:   uuid.New(),
		Name: "test",
		Request: packaging.HttpRequest{
			Method: "GET",
			URL:    "http://localhost:8888/account?id=$X$",
			Headers: []packaging.KeyValue{
				{
					Key:   "Host",
					Value: "localhost:8888",
				},
			},
		},
		Input: Link{
			To: LinkDirection{
				Node:      fuzzer.ID(),
				Connector: fuzzer.GetInput().ID,
			},
		},
		Output: nOutput,
		Error:  nError,
		Nodes: []Node{
			fuzzer,
			verifier,
		},
		Links: []Link{
			{
				From: LinkDirection{
					Node:      fuzzer.ID(),
					Connector: fuzzer.GetOutputs()[0].ID,
				},
				To: LinkDirection{
					Node:      verifier.ID(),
					Connector: verifier.GetInput().ID,
				},
			},
			{
				From: LinkDirection{
					Node:      verifier.ID(),
					Connector: verifierGood.ID,
				},
				To: LinkDirection{
					Node:      nOutput.ID(),
					Connector: nOutput.GetInput().ID,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	updates := make(chan Update)
	defer close(updates)

	go func() {
		for update := range updates {
			fmt.Printf("update: %s: %s\n", update.Node, update.Message)
		}
	}()

	require.NoError(t, workflow.Run(ctx, updates))

	assert.True(t, strings.HasSuffix(strings.Split(output.String(), "\n")[0], fmt.Sprintf("($X$=%d)", secret)))

}
