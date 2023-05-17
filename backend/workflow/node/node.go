package node

import (
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Channel string

const maxThreadsPerNode = 10

const (
	ChannelStdout   Channel = "stdout"
	ChannelStderr   Channel = "stderr"
	ChannelActivity Channel = "activity"
)

type Output struct {
	Node    uuid.UUID
	Channel Channel
	Message string
}

type OutputM struct {
	Node    string `json:"node"`
	Channel string `json:"channel"`
	Message string `json:"message"`
}

func (o Output) Pack() OutputM {
	return OutputM{
		Node:    o.Node.String(),
		Channel: string(o.Channel),
		Message: o.Message,
	}
}

type Type int

const (
	TypeUnknown Type = iota
	TypeFuzzer
	TypeOutput
	TypeStatusFilter
	TypeRequest
	TypeStart
	TypeSender
	TypeVariables
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
	Run(context.Context, map[string]transmission.Transmission, chan<- Output, bool) (<-chan OutputInstance, <-chan error)
	Validate(params map[string]transmission.Transmission) error
}

type OutputInstance struct {
	OutputName string
	Current    int
	Total      int
	Complete   bool
	Data       transmission.Transmission
}
