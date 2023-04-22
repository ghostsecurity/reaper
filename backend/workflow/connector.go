package workflow

import "github.com/google/uuid"

type Connector struct {
	ID   uuid.UUID
	Name string
	Type TransmissionType
}

func NewConnector(name string, t TransmissionType) Connector {
	return Connector{
		ID:   uuid.New(),
		Name: name,
		Type: t,
	}
}

type Connectors []Connector

func (c Connectors) Find(id uuid.UUID) (Connector, bool) {
	for _, conn := range c {
		if conn.ID == id {
			return conn, true
		}
	}
	return Connector{}, false
}

func (c Connectors) FindByName(name string) (Connector, bool) {
	for _, conn := range c {
		if conn.Name == name {
			return conn, true
		}
	}
	return Connector{}, false
}
