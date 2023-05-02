package node

import (
	"strings"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
)

type Connector struct {
	Name        string                  `json:"name"`
	Type        transmission.ParentType `json:"type"`
	Linkable    bool                    `json:"linkable"`
	Description string                  `json:"description"`
}

func NewConnector(name string, t transmission.ParentType, linkable bool, description ...string) Connector {
	return Connector{
		Name:        name,
		Type:        t,
		Linkable:    linkable,
		Description: strings.Join(description, ", "),
	}
}

type Connectors []Connector

func (c Connectors) FindByName(name string) (Connector, bool) {
	for _, conn := range c {
		if conn.Name == name {
			return conn, true
		}
	}
	return Connector{}, false
}
