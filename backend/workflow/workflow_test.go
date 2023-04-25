package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
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

	fuzzer := node.NewFuzzer()
	fuzzer.SetStaticInputValues(map[string]transmission.Transmission{
		"placeholder": transmission.NewString("$ID$"),
		"list":        transmission.NewNumericRangeIterator(0, 100),
	})
	verifier := node.NewVerifier()
	verifier.SetStaticInputValues(map[string]transmission.Transmission{
		"min": transmission.NewInt(200),
		"max": transmission.NewInt(200),
	})

	nOutput := node.NewOutput()
	nError := node.NewOutput()

	// start workflow with input node and output node
	workflow := &Workflow{
		ID:   uuid.New(),
		Name: "test",
		Request: packaging.HttpRequest{
			Method: "GET",
			URL:    "http://localhost:8888/account?id=$ID$",
			Headers: []packaging.KeyValue{
				{
					Key:   "Host",
					Value: "localhost:8888",
				},
			},
		},
		Input: node.Link{
			To: node.LinkDirection{
				Node:      fuzzer.ID(),
				Connector: "request",
			},
		},
		Output: nOutput,
		Error:  nError,
		Nodes: []node.Node{
			fuzzer,
			verifier,
		},
		Links: []node.Link{
			{
				From: node.LinkDirection{
					Node:      fuzzer.ID(),
					Connector: "output",
				},
				To: node.LinkDirection{
					Node:      verifier.ID(),
					Connector: "response",
				},
			},
			{
				From: node.LinkDirection{
					Node:      verifier.ID(),
					Connector: "good",
				},
				To: node.LinkDirection{
					Node:      nOutput.ID(),
					Connector: "input",
				},
			},
		},
	}

	require.NoError(t, workflow.Validate())

	t.Run("run", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		updates := make(chan Update)
		defer close(updates)

		go func() {
			for update := range updates {
				_ = update
				//fmt.Printf("update: %s: %s\n", update.Node, update.Message)
			}
		}()

		require.NoError(t, workflow.Run(ctx, updates))

		assert.True(t, strings.HasSuffix(strings.Split(nOutput.String(), "\n")[0], fmt.Sprintf("($ID$=%d)", secret)))

	})

	t.Run("save to disk, reload and run", func(t *testing.T) {
		data, err := json.Marshal(workflow)
		require.NoError(t, err)

		fmt.Println(string(data))

		var w Workflow
		require.NoError(t, json.Unmarshal(data, &w))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		updates := make(chan Update)
		defer close(updates)

		go func() {
			for update := range updates {
				_ = update
				//fmt.Printf("update: %s: %s\n", update.Node, update.Message)
			}
		}()

		require.NoError(t, w.Run(ctx, updates))

		assert.True(t, strings.HasSuffix(strings.Split(nOutput.String(), "\n")[0], fmt.Sprintf("($ID$=%d)", secret)))
	})
}
