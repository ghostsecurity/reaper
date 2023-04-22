package workflow

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Node interface {
	ID() uuid.UUID
	Name() string
	GetInput() Connector
	GetOutputs() Connectors
	Run(context.Context, map[uuid.UUID]Transmission) (<-chan OutputInstance, <-chan error)
	// TODO Add json marshal/unmarshal
}

type OutputInstance struct {
	OutputID uuid.UUID
	Current  int
	Total    int
	Complete bool
	Data     Transmission
}
