package node

import (
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type NodeType int

const (
	TypeUnknown NodeType = iota
	TypeFuzzer
	TypeOutput
	TypeVerifier
)

type Node interface {
	ID() uuid.UUID
	SetID(uuid.UUID)
	Type() NodeType
	Name() string
	GetInputs() Connectors
	SetStaticInputValues(map[string]transmission.Transmission)
	GetOutputs() Connectors
	GetVars() *VarStorage
	SetVars(*VarStorage)
	Run(context.Context, map[string]transmission.Transmission) (<-chan OutputInstance, <-chan error)
}

type OutputInstance struct {
	OutputName string
	Current    int
	Total      int
	Complete   bool
	Data       transmission.Transmission
}
