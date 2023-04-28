package node

import (
	"io"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Type int

const (
	TypeUnknown Type = iota
	TypeFuzzer
	TypeOutput
	TypeStatusFilter
	TypeRequest
	TypeStart
	TypeSender
)

type Node interface {
	IsReadOnly() bool
	ID() uuid.UUID
	SetID(uuid.UUID)
	SetName(string)
	Type() Type
	Name() string
	GetInputs() Connectors
	SetStaticInputValues(map[string]transmission.Transmission) error
	AddStaticInputValue(string, transmission.Transmission) error
	GetInjections() map[string]transmission.Transmission
	GetOutputs() Connectors
	GetVars() *VarStorage
	SetVars(*VarStorage)
	Run(context.Context, map[string]transmission.Transmission, io.Writer, io.Writer) (<-chan OutputInstance, <-chan error)
	Validate(params map[string]transmission.Transmission) error
}

type OutputInstance struct {
	OutputName string
	Current    int
	Total      int
	Complete   bool
	Data       transmission.Transmission
}
