package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
	"github.com/google/uuid"
)

type workflowM struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Request     packaging.HttpRequest  `json:"request"`
	Input       node.Link              `json:"input"`
	Output      nodeM                  `json:"output"`
	Error       nodeM                  `json:"error"`
	Nodes       []nodeM                `json:"nodes"`
	Links       []node.Link            `json:"links"`
	Positioning map[uuid.UUID]Position `json:"positioning"`
}

func (w *Workflow) MarshalJSON() ([]byte, error) {

	mOutput, err := toNodeM(w.Output)
	if err != nil {
		return nil, err
	}

	mError, err := toNodeM(w.Error)
	if err != nil {
		return nil, err
	}

	mNodes := make([]nodeM, len(w.Nodes))
	for i, node := range w.Nodes {
		nm, err := toNodeM(node)
		if err != nil {
			return nil, err
		}
		mNodes[i] = *nm
	}

	m := workflowM{
		ID:          w.ID,
		Name:        w.Name,
		Request:     w.Request,
		Input:       w.Input,
		Output:      *mOutput,
		Error:       *mError,
		Nodes:       mNodes,
		Links:       w.Links,
		Positioning: w.Positioning,
	}

	return json.Marshal(m)
}

func (w *Workflow) UnmarshalJSON(data []byte) error {
	var m workflowM
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*w = Workflow{
		ID:          m.ID,
		Name:        m.Name,
		Request:     m.Request,
		Input:       m.Input,
		Nodes:       make([]node.Node, len(m.Nodes)),
		Links:       m.Links,
		Positioning: m.Positioning,
	}

	mOut, err := m.Output.ToNode()
	if err != nil {
		return err
	}
	w.Output = mOut
	mErr, err := m.Error.ToNode()
	if err != nil {
		return err
	}
	w.Error = mErr

	for i, node := range m.Nodes {
		n, err := node.ToNode()
		if err != nil {
			return err
		}
		w.Nodes[i] = n
	}
	return nil
}

type nodeM struct {
	Id   uuid.UUID        `json:"id"`
	Type node.NodeType    `json:"type"`
	Vars *node.VarStorage `json:"vars"`
}

func toNodeM(n node.Node) (*nodeM, error) {
	return &nodeM{
		Id:   n.ID(),
		Type: n.Type(),
		Vars: n.GetVars(),
	}, nil
}

func (n *nodeM) ToNode() (node.Node, error) {
	var real node.Node
	switch n.Type {
	case node.TypeFuzzer:
		real = node.NewFuzzer()
	case node.TypeVerifier:
		real = node.NewVerifier()
	case node.TypeOutput:
		real = node.NewOutput()
	default:
		return nil, fmt.Errorf("unknown node type: %v", n.Type)
	}
	real.SetID(n.Id)
	real.SetVars(n.Vars)
	return real, nil
}
