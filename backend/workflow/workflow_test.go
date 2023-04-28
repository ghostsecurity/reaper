package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
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
	require.NoError(t, fuzzer.SetStaticInputValues(map[string]transmission.Transmission{
		"placeholder": transmission.NewString("$ID$"),
		"list":        transmission.NewNumericRangeIterator(0, 100),
	}))
	verifier := node.NewStatusFilter()
	require.NoError(t, verifier.SetStaticInputValues(map[string]transmission.Transmission{
		"min": transmission.NewInt(200),
		"max": transmission.NewInt(200),
	}))

	nOutput := node.NewOutput()
	nError := node.NewOutput()
	nError.SetStaticInputValues(map[string]transmission.Transmission{
		"stdout": transmission.NewBoolean(false),
		"stderr": transmission.NewBoolean(true),
	})

	reqNode := node.NewRequest()
	reqNode.SetStaticInputValues(map[string]transmission.Transmission{
		"input": transmission.NewRequest(packaging.HttpRequest{
			Method: "GET",
			URL:    "http://localhost:8888/account?id=$ID$",
			Headers: []packaging.KeyValue{
				{
					Key:   "Host",
					Value: "localhost:8888",
				},
			},
		}),
	})

	sender := node.NewSender()

	flow := NewWorkflow()
	start := flow.Nodes[0]
	flow.Nodes = []node.Node{
		start,
		reqNode,
		fuzzer,
		verifier,
		nOutput,
		nError,
		sender,
	}

	flow.Links = []node.Link{
		{
			From: node.LinkDirection{
				Node:      start.ID(),
				Connector: "output",
			},
			To: node.LinkDirection{
				Node:      fuzzer.ID(),
				Connector: "start",
			},
		},
		{
			From: node.LinkDirection{
				Node:      reqNode.ID(),
				Connector: "output",
			},
			To: node.LinkDirection{
				Node:      sender.ID(),
				Connector: "request",
			},
		},
		{
			From: node.LinkDirection{
				Node:      fuzzer.ID(),
				Connector: "output",
			},
			To: node.LinkDirection{
				Node:      sender.ID(),
				Connector: "replacements",
			},
		},
		{
			From: node.LinkDirection{
				Node:      sender.ID(),
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
	}

	require.NoError(t, flow.Validate())

	t.Run("run", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		updates := make(chan Update)
		defer close(updates)

		go func() {
			for update := range updates {
				fmt.Printf("update: %s: %s\n", update.Node, update.Message)
			}
		}()

		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)
		require.NoError(t, flow.Run(ctx, updates, stdout, stderr))
		assert.True(t, strings.HasSuffix(strings.Split(stdout.String(), "\n")[0], fmt.Sprintf("[$ID$=%d]", secret)))
	})

	t.Run("save to disk, reload and run", func(t *testing.T) {

		data, err := json.Marshal(flow)
		require.NoError(t, err)

		var w Workflow
		require.NoError(t, json.Unmarshal(data, &w))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		updates := make(chan Update)
		defer close(updates)

		go func() {
			for update := range updates {
				fmt.Printf("update: %s: %s\n", update.Node, update.Message)
			}
		}()

		stdout := bytes.NewBuffer(nil)
		stderr := io.Discard
		require.NoError(t, w.Run(ctx, updates, stdout, stderr))
		assert.True(t, strings.HasSuffix(strings.Split(stdout.String(), "\n")[0], fmt.Sprintf("[$ID$=%d]", secret)))
	})
}
