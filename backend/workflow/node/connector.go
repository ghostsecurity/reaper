package node

import "github.com/ghostsecurity/reaper/backend/workflow/transmission"

type Connector struct {
	Name string                  `json:"name"`
	Type transmission.ParentType `json:"type"`
}

func NewConnector(name string, t transmission.ParentType) Connector {
	return Connector{
		Name: name,
		Type: t,
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
