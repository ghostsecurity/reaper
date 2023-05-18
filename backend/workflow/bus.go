package workflow

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"

	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/ghostsecurity/reaper/backend/workflow/node"
)

type Bus struct {
	start      node.Node
	nodes      []node.Node
	links      []node.Link
	inputs     map[uuid.UUID]chan node.Input
	routes     map[string][]Route
	statuses   map[uuid.UUID]Update
	updateChan chan<- Update
	mu         sync.Mutex
	statusMu   sync.Mutex
	inputsMu   sync.RWMutex
	complete   map[uuid.UUID]bool
	completeMu sync.Mutex
	aborted    bool
	abortedMu  sync.RWMutex
}

type Route struct {
	inputName string
	node      uuid.UUID
}

func NewBus(start node.Node, sc chan<- Update) *Bus {
	return &Bus{
		start:      start,
		updateChan: sc,
		inputs:     make(map[uuid.UUID]chan node.Input),
		routes:     make(map[string][]Route),
		statuses:   make(map[uuid.UUID]Update),
		complete:   make(map[uuid.UUID]bool),
	}
}

func (b *Bus) AddNode(n node.Node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}
	for _, existing := range b.nodes {
		if existing.ID() == n.ID() {
			return fmt.Errorf("node with id %s already exists", n.ID())
		}
	}
	b.nodes = append(b.nodes, n)
	return nil
}

func (b *Bus) AddLink(l node.Link) error {
	for _, existing := range b.links {
		if existing == l {
			return fmt.Errorf("link already exists")
		}
	}
	b.links = append(b.links, l)
	return nil
}

func (b *Bus) buildRoutes() {
	for _, l := range b.links {
		key := l.From.Node.String() + ":" + l.From.Connector
		list := b.routes[key]
		list = append(list, Route{
			inputName: l.To.Connector,
			node:      l.To.Node,
		})
		b.routes[key] = list
	}
}

// closeNodeInput is called when a node has no further input.
func (b *Bus) closeNodeInput(n node.Node) {
	b.inputsMu.Lock()
	defer b.inputsMu.Unlock()
	if in, ok := b.inputs[n.ID()]; ok {
		close(in)
		delete(b.inputs, n.ID())
	}
}

func (b *Bus) abort() {
	b.abortedMu.Lock()
	b.aborted = true
	b.abortedMu.Unlock()
}

func (b *Bus) isAborted() bool {
	b.abortedMu.RLock()
	defer b.abortedMu.RUnlock()
	return b.aborted
}

// closeNodeOutput is called when a node has no further output.
func (b *Bus) closeNodeOutput(n node.Node) {

	b.closeNodeInput(n)

	for _, other := range b.nodes {
		if other.ID() == n.ID() {
			continue
		}
		b.inputsMu.RLock()
		_, remaining := b.inputs[other.ID()]
		b.inputsMu.RUnlock()
		if !remaining {
			continue
		}
		var countDirect int
		var countOther int
		for _, l := range b.links {
			if l.To.Node == other.ID() {
				if l.From.Node == n.ID() {
					countDirect++
				} else {
					b.inputsMu.RLock()
					_, incomplete := b.inputs[l.From.Node]
					b.inputsMu.RUnlock()
					if incomplete {
						countOther++
					}
				}
			}
		}
		if countDirect > 0 && countOther == 0 {
			b.closeNodeInput(other)
		}
	}
}

func (b *Bus) Run(ctx context.Context, output chan<- node.Output) error {

	b.mu.Lock()
	defer b.mu.Unlock()
	b.inputs = make(map[uuid.UUID]chan node.Input)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	var firstNodeError error
	var errMu sync.Mutex

	for _, n := range b.nodes {
		b.inputs[n.ID()] = make(chan node.Input)
	}

	b.buildRoutes()

	for _, n := range b.nodes {

		var hasInputConnected bool
		for _, l := range b.links {
			if l.To.Node == n.ID() {
				hasInputConnected = true
				break
			}
		}
		if !hasInputConnected && b.start.ID() != n.ID() {
			b.updateStatus(ctx, Update{
				Node:    n.ID(),
				Name:    n.Name(),
				Status:  NodeStatusSuccess,
				Message: "Success (no inputs)",
			})
			b.closeNodeInput(n)
			continue
		}

		b.inputsMu.RLock()
		in := b.inputs[n.ID()]
		b.inputsMu.RUnlock()
		out := make(chan node.OutputInstance, 100)
		wg.Add(2)
		go func(n node.Node, in chan node.Input, out chan node.OutputInstance) {
			defer wg.Done()
			if len(n.GetInjections()) > 0 {
				b.updateStatus(ctx, Update{
					Node:    n.ID(),
					Name:    n.Name(),
					Status:  NodeStatusRunning,
					Message: "Input(s) injected...",
				})
			} else {
				b.updateStatus(ctx, Update{
					Node:    n.ID(),
					Name:    n.Name(),
					Status:  NodeStatusPending,
					Message: "Waiting for input(s)...",
				})
			}
			defer close(out)
			if err := n.Start(ctx, in, out, output); err != nil {
				if errors.Is(err, context.Canceled) {
					if !b.isAborted() {
						b.abort()
					}
					b.updateStatus(ctx, Update{
						Node:    n.ID(),
						Name:    n.Name(),
						Status:  NodeStatusAborted,
						Message: "Aborted",
					})
				} else {
					if !b.isAborted() {
						b.abort()
					}
					b.updateStatus(ctx, Update{
						Node:    n.ID(),
						Name:    n.Name(),
						Status:  NodeStatusError,
						Message: "Error: " + err.Error(),
					})
					cancel()
				}
				errMu.Lock()
				if firstNodeError == nil {
					firstNodeError = err
				}
				errMu.Unlock()
			}

		}(n, in, out)
		go func(n node.Node, out chan node.OutputInstance) {
			defer wg.Done()
			defer b.checkDeadlock(n.ID())
			defer b.closeNodeOutput(n)
			for {
				msg, ok := <-out
				if !ok {
					b.updateStatus(ctx, Update{
						Node:    n.ID(),
						Name:    n.Name(),
						Status:  NodeStatusSuccess,
						Message: "Success (routine complete)",
					})
					return
				}

				select {
				case <-ctx.Done():
					return
				default:
				}

				routes := b.routes[n.ID().String()+":"+msg.OutputName]

				for _, route := range routes {
					b.updateStatus(ctx, Update{
						Node:    route.node,
						Status:  NodeStatusRunning,
						Message: "Running...",
					})
					b.inputsMu.RLock()
					c, ok := b.inputs[route.node]
					b.inputsMu.RUnlock()
					if !ok {
						continue
					}
					if b.isAborted() {
						return
					}
					func() {
						// last ditch to catch rogue chan writes
						defer func() {
							if p := recover(); p != nil {
								fmt.Println("failed to write to input channel:", p)
							}
						}()
						select {
						case <-ctx.Done():
							break
						case c <- node.Input{
							Last: msg.Complete,
							Data: map[string]transmission.Transmission{
								route.inputName: msg.Data,
							},
						}:
						}
					}()
				}

				b.checkDeadlock(n.ID())
			}
		}(n, out)
	}

	b.inputsMu.RLock()
	startInput, ok := b.inputs[b.start.ID()]
	if !ok {
		b.inputsMu.RUnlock()
		cancel()
		wg.Wait()
		return fmt.Errorf("start node not found")
	}

	b.updateStatus(ctx, Update{
		Node:    b.start.ID(),
		Name:    b.start.Name(),
		Status:  NodeStatusRunning,
		Message: "Running...",
	})
	select {
	case <-ctx.Done():
		b.inputsMu.RUnlock()
		wg.Wait()
		return ctx.Err()
	case startInput <- node.Input{
		Last: true,
		Data: nil,
	}:
		b.inputsMu.RUnlock()
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
SAFETY:
	for {
		select {
		case <-func() chan struct{} {
			c := make(chan struct{})
			go func() {
				defer close(c)
				wg.Wait()
			}()
			return c
		}():
			break SAFETY
		case <-ticker.C:
			var remaining []uuid.UUID
			b.inputsMu.RLock()
			for n := range b.inputs {
				remaining = append(remaining, n)
			}
			b.inputsMu.RUnlock()
			for _, r := range remaining {
				b.checkDeadlock(r)
			}
		}
	}

	for _, n := range b.nodes {
		if old, ok := b.statuses[n.ID()]; !ok || (old.Status != NodeStatusSuccess && old.Status != NodeStatusError && old.Status != NodeStatusAborted) {
			if firstNodeError == nil {
				b.updateStatus(ctx, Update{
					Node:    n.ID(),
					Name:    n.Name(),
					Status:  NodeStatusSuccess,
					Message: "Success (defaulted)",
				})
			} else {
				b.updateStatus(ctx, Update{
					Node:    n.ID(),
					Name:    n.Name(),
					Status:  NodeStatusAborted,
					Message: "Aborted",
				})
			}
		}
	}

	return firstNodeError
}

func (b *Bus) isNodeDead(id uuid.UUID) bool {
	b.statusMu.Lock()
	defer b.statusMu.Unlock()
	if status, ok := b.statuses[id]; ok {
		return status.Status == NodeStatusAborted || status.Status == NodeStatusError
	}
	return false
}

func (b *Bus) updateStatus(_ context.Context, update Update) {
	b.statusMu.Lock()
	defer b.statusMu.Unlock()
	if old, ok := b.statuses[update.Node]; ok {
		if old.Status == update.Status && old.Message == update.Message {
			return
		}
		if old.Status == NodeStatusSuccess || old.Status == NodeStatusError || old.Status == NodeStatusAborted {
			b.statuses[update.Node] = update
			return
		}
	}
	b.statuses[update.Node] = update
	b.updateChan <- update
}

func (b *Bus) checkDeadlock(fromOutputNode uuid.UUID) {

	b.inputsMu.RLock()
	defer b.inputsMu.RUnlock()

	// first let's grab our remaining nodes which are running
	var remainingNodes []uuid.UUID
	for _, n := range b.nodes {
		var linked bool
		for _, link := range b.links {
			if link.From.Node == fromOutputNode && link.To.Node == n.ID() {
				linked = true
				break
			}
		}
		if !linked {
			continue
		}
		if _, ok := b.inputs[n.ID()]; ok {
			remainingNodes = append(remainingNodes, n.ID())
		}
	}

	if len(remainingNodes) == 0 {
		return
	}

	// we need to look for circular dependencies which are strangled from all other inputs
	for _, n := range remainingNodes {
		chained, circ := b.getChainedInputNodes(n, nil)
		if !circ || len(chained) == 0 {
			continue
		}
		// we have a circular dependency, let's see if it's strangled
		strangled := true
		for _, chain := range chained {
			if !b.isCircular(chain) {
				if _, ok := b.inputs[chain]; ok {
					// this chain is not strangled
					strangled = false
					break
				}
			}
		}
		if strangled {
			b.inputsMu.RUnlock()
			b.shutDownChain(chained)
			b.inputsMu.RLock()
			return
		}
	}

}

func (b *Bus) shutDownChain(chain []uuid.UUID) {

	var nodes []node.Node

	// check if anything in the chain is busy - if not, we can probably shut it down
	for _, id := range chain {
		var n node.Node
		for _, no := range b.nodes {
			if no.ID() == id {
				n = no
				break
			}
		}
		if n == nil {
			continue
		}
		if n.Busy() {
			return
		}
		nodes = append(nodes, n)
	}

	for _, n := range nodes {
		b.closeNodeInput(n)
	}
}

func (b *Bus) isCircular(n uuid.UUID) bool {
	_, circular := b.getChainedInputNodes(n, nil)
	return circular
}

// returns all input nodes to a given node, and all input nodes to those nodes, and so on
func (b *Bus) getChainedInputNodes(from uuid.UUID, used []uuid.UUID) ([]uuid.UUID, bool) {
	nodes := []uuid.UUID{
		from,
	}
	used = append(used, from)
	var circular bool
	for _, link := range b.links {
		if link.To.Node == from {
			// node has completed, not interested
			if _, ok := b.inputs[link.From.Node]; !ok {
				continue
			}
			var exists bool
			for _, u := range used {
				if u == link.From.Node {
					exists = true
					break
				}
			}
			if exists {
				circular = true
				continue
			}
			chained, circ := b.getChainedInputNodes(link.From.Node, used)
			circular = circular || circ
			nodes = append(nodes, chained...)
		}
	}
	return nodes, circular
}
