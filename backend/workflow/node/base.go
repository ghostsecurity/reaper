package node

import (
	"sync"
	"time"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type noInjections struct{}

func (n noInjections) GetInjections() map[string]transmission.Transmission {
	return nil
}

type base struct {
	*VarStorage
	id       uuid.UUID
	name     string
	t        Type
	readonly bool
	busy     bool
	last     time.Time
	busyMu   sync.RWMutex
}

func newBase(name string, t Type, readonly bool, vars *VarStorage) *base {
	return &base{
		VarStorage: vars,
		id:         uuid.New(),
		name:       name,
		t:          t,
		readonly:   readonly,
	}
}

func (b *base) Busy() bool {
	b.busyMu.RLock()
	defer b.busyMu.RUnlock()
	return b.busy
}

func (b *base) LastInput() time.Time {
	b.busyMu.RLock()
	defer b.busyMu.RUnlock()
	return b.last
}

func (b *base) setBusy(busy bool) {
	b.busyMu.Lock()
	defer b.busyMu.Unlock()
	b.busy = busy
	b.last = time.Now()
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

func (b *base) MergeVars(vars *VarStorage) {
	b.VarStorage.Merge(vars)
}

func (b *base) tryOut(ctx context.Context, out chan<- OutputInstance, instance OutputInstance) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- instance:
		return nil
	}
}
