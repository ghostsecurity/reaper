package workflow

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/ghostsecurity/reaper/backend/workflow/node"
)

/*
	NOTE:
	A lot of the strangeness in here existed to help wails create js bindings - we can start to simplify this now...
*/

type WorkflowM struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Nodes       []NodeM             `json:"nodes"`
	Links       []LinkM             `json:"links"`
	Positioning map[string]Position `json:"positioning"`
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

	mNodes := make([]NodeM, len(w.Nodes))
	for i, node := range w.Nodes {
		nm, err := ToNodeM(node)
		if err != nil {
			return nil, err
		}
		mNodes[i] = *nm
	}

	return &WorkflowM{
		ID:          w.ID.String(),
		Name:        w.Name,
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
		Nodes:       make([]node.Node, len(m.Nodes)),
		Links:       fromLinkMs(m.Links),
		Positioning: fromPositioningM(m.Positioning),
	}

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
	Id       string            `json:"id"`
	Name     string            `json:"name"`
	Type     int               `json:"type"`
	Vars     *node.VarStorageM `json:"vars"`
	ReadOnly bool              `json:"readonly"`
}

func ToNodeM(n node.Node) (*NodeM, error) {
	packed, err := n.GetVars().Pack()
	if err != nil {
		return nil, err
	}
	return &NodeM{
		Id:       n.ID().String(),
		Name:     n.Name(),
		Type:     int(n.Type()),
		Vars:     packed,
		ReadOnly: n.IsReadOnly(),
	}, nil
}

func (n *NodeM) ToNode() (node.Node, error) {
	real, err := node.FromType(node.Type(n.Type))
	if err != nil {
		return nil, err
	}
	real.SetID(toUUIDOrNil(n.Id))
	real.SetName(n.Name)
	unpacked, err := n.Vars.Unpack()
	if err != nil {
		return nil, err
	}
	real.MergeVars(unpacked)
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

type VarsM struct {
	Values map[string]interface{} `json:"static"`
}
