package node

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
)

type Channel string

const (
	ChannelStdout   Channel = "stdout"
	ChannelStderr   Channel = "stderr"
	ChannelActivity Channel = "activity"
)

type Input struct {
	Last bool
	Data map[string]transmission.Transmission
}

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
	TypeDelay
	TypeExtractor
	TypeIf
)

func FromType(t Type) (Node, error) {
	var real Node
	switch t {
	case TypeFuzzer:
		real = NewFuzzer()
	case TypeStatusFilter:
		real = NewStatusFilter()
	case TypeOutput:
		real = NewOutput()
	case TypeRequest:
		real = NewRequest()
	case TypeStart:
		real = NewStart()
	case TypeSender:
		real = NewSender()
	case TypeVariables:
		real = NewVars()
	case TypeDelay:
		real = NewDelay()
	case TypeExtractor:
		real = NewExtractor()
	case TypeIf:
		real = NewIf()
	default:
		return nil, fmt.Errorf("unknown node type: %v", t)
	}
	return real, nil
}

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
	GetOutputs() Connectors
	GetVars() *VarStorage
	SetVars(*VarStorage)
	MergeVars(*VarStorage)
	Validate(params map[string]transmission.Transmission) error
	LastInput() time.Time
	Busy() bool

	GetInjections() map[string]transmission.Transmission
	Start(context.Context, <-chan Input, chan<- OutputInstance, chan<- Output) error
}

type OutputInstance struct {
	OutputName string
	Current    int
	Total      int
	Complete   bool
	Data       transmission.Transmission
}
