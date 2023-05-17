package node

import (
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
)

type noInjections struct{}

func (n noInjections) GetInjections() map[string]transmission.Transmission {
	return nil
}

type base struct {
	*VarStorage
	id       uuid.UUID
	name     string
	input    chan map[string]transmission.Transmission
	output   chan OutputInstance
	t        Type
	readonly bool
}

func newBase(name string, t Type, readonly bool, vars *VarStorage) *base {
	return &base{
		VarStorage: vars,
		id:         uuid.New(),
		name:       name,
		input:      make(chan map[string]transmission.Transmission, maxThreadsPerNode),
		output:     make(chan OutputInstance),
		t:          t,
		readonly:   readonly,
	}
}

func (b *base) ID() uuid.UUID {
	return b.id
}

func (b *base) SetID(id uuid.UUID) {
	b.id = id
}

func (b *base) Name() string {
	return b.name
}

func (b *base) SetName(name string) {
	b.name = name
}

func (b *base) Type() Type {
	return b.t
}

func (b *base) IsReadOnly() bool {
	return b.readonly
}

func (b *base) GetVars() *VarStorage {
	return b.VarStorage
}

func (b *base) SetVars(vars *VarStorage) {
	b.VarStorage = vars
}
