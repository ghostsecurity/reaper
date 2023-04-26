package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
	"github.com/google/uuid"
)

/*
	NOTE:
	A lot of the strangeness in here exists to help wails create js bindings.
*/

type WorkflowM struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Request     packaging.HttpRequest `json:"request"`
	Input       LinkM                 `json:"input"`
	Output      NodeM                 `json:"output"`
	Error       NodeM                 `json:"error"`
	Nodes       []NodeM               `json:"nodes"`
	Links       []LinkM               `json:"links"`
	Positioning map[string]Position   `json:"positioning"`
}

type LinkDirectionM struct {
	Node      string `json:"node"`
	Connector string `json:"connector"`
}

type LinkM struct {
	From       LinkDirectionM `json:"from"`
	To         LinkDirectionM `json:"to"`
	Annotation string         `json:"annotation"`
}

func (w *Workflow) Pack() (*WorkflowM, error) {
	mOutput, err := toNodeM(w.Output)
	if err != nil {
		return nil, err
	}

	mError, err := toNodeM(w.Error)
	if err != nil {
		return nil, err
	}

	mNodes := make([]NodeM, len(w.Nodes))
	for i, node := range w.Nodes {
		nm, err := toNodeM(node)
		if err != nil {
			return nil, err
		}
		mNodes[i] = *nm
	}

	return &WorkflowM{
		ID:          w.ID.String(),
		Name:        w.Name,
		Request:     w.Request,
		Input:       toLinkM(w.Input),
		Output:      *mOutput,
		Error:       *mError,
		Nodes:       mNodes,
		Links:       toLinkMs(w.Links),
		Positioning: toPositioningM(w.Positioning),
	}, nil
}

func (w *Workflow) MarshalJSON() ([]byte, error) {
	m, err := w.Pack()
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func toPositioningM(p map[uuid.UUID]Position) map[string]Position {
	output := make(map[string]Position)
	for k, v := range p {
		output[k.String()] = v
	}
	return output
}

func fromPositioningM(p map[string]Position) map[uuid.UUID]Position {
	output := make(map[uuid.UUID]Position)
	for k, v := range p {
		output[toUUIDOrNil(k)] = v
	}
	return output
}

func toUUIDOrNil(u string) uuid.UUID {
	if p, err := uuid.Parse(u); err == nil {
		return p
	}
	return uuid.Nil
}

func (m *WorkflowM) Unpack() (*Workflow, error) {
	w := Workflow{
		ID:          toUUIDOrNil(m.ID),
		Name:        m.Name,
		Request:     m.Request,
		Input:       fromLinkM(m.Input),
		Nodes:       make([]node.Node, len(m.Nodes)),
		Links:       fromLinkMs(m.Links),
		Positioning: fromPositioningM(m.Positioning),
	}

	mOut, err := m.Output.ToNode()
	if err != nil {
		return nil, err
	}
	w.Output = mOut
	mErr, err := m.Error.ToNode()
	if err != nil {
		return nil, err
	}
	w.Error = mErr

	for i, node := range m.Nodes {
		n, err := node.ToNode()
		if err != nil {
			return nil, err
		}
		w.Nodes[i] = n
	}
	return &w, nil
}

func (w *Workflow) UnmarshalJSON(data []byte) error {
	var m WorkflowM
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	ww, err := m.Unpack()
	if err != nil {
		return err
	}
	*w = *ww
	return nil
}

type NodeM struct {
	Id   uuid.UUID        `json:"id"`
	Type node.NodeType    `json:"type"`
	Vars *node.VarStorage `json:"vars"`
}

func toNodeM(n node.Node) (*NodeM, error) {
	return &NodeM{
		Id:   n.ID(),
		Type: n.Type(),
		Vars: n.GetVars(),
	}, nil
}

func (n *NodeM) ToNode() (node.Node, error) {
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

func toLinkM(l node.Link) LinkM {
	return LinkM{
		From:       toLinkDirectionM(l.From),
		To:         toLinkDirectionM(l.To),
		Annotation: l.Annotation,
	}
}

func toLinkMs(links []node.Link) []LinkM {
	mLinks := make([]LinkM, len(links))
	for i, link := range links {
		mLinks[i] = toLinkM(link)
	}
	return mLinks
}

func toLinkDirectionM(ld node.LinkDirection) LinkDirectionM {
	return LinkDirectionM{
		Node:      ld.Node.String(),
		Connector: ld.Connector,
	}
}

func fromLinkM(l LinkM) node.Link {
	return node.Link{
		From:       fromLinkDirectionM(l.From),
		To:         fromLinkDirectionM(l.To),
		Annotation: l.Annotation,
	}
}

func fromLinkMs(links []LinkM) []node.Link {
	mLinks := make([]node.Link, len(links))
	for i, link := range links {
		mLinks[i] = fromLinkM(link)
	}
	return mLinks
}

func fromLinkDirectionM(ld LinkDirectionM) node.LinkDirection {
	return node.LinkDirection{
		Node:      toUUIDOrNil(ld.Node),
		Connector: ld.Connector,
	}
}
