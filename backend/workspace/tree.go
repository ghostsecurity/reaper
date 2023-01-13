package workspace

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Tree struct {
	Root StructureNode `json:"root"`
}

func (t *Tree) Update(request *http.Request) bool {
	return t.Root.Update(append([]string{
		request.URL.Hostname(),
	},
		strings.Split(request.URL.Path, "/")...,
	))
}

func (t *Tree) Structure() []StructureNode {
	return t.Root.Children
}

type StructureNode struct {
	Name     string          `json:"name"`
	Children []StructureNode `json:"children"`
}

func (t *StructureNode) Update(parts []string) bool {
	var filtered []string
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	if len(filtered) == 0 {
		return false
	}
	for i, node := range t.Children {
		if node.Name == filtered[0] {
			return t.Children[i].Update(filtered[1:])
		}
	}
	hostNode := StructureNode{
		Name: filtered[0],
	}
	_ = hostNode.Update(filtered[1:])
	t.Children = append(t.Children, hostNode)
	return true
}

func (t *StructureNode) MarshalJSON() ([]byte, error) {
	if t.Children == nil {
		t.Children = []StructureNode{}
	}
	return json.Marshal(map[string]interface{}{
		"name":     t.Name,
		"children": t.Children,
	})
}
