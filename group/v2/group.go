package v2

import (
	"fmt"
	"sync"
)

type CreaterFunc[T any] func(string) T

// Group is a lazy load container.
type Group[T any] struct {
	new  CreaterFunc[T]
	vals map[string]T
	sync.RWMutex
}

// NewGroup news a group container.
func NewGroup[T any](new CreaterFunc[T]) *Group[T] {
	if new == nil {
		panic(fmt.Errorf("container.group: can't assign a nil to the new function"))
	}
	return &Group[T]{
		new:  new,
		vals: make(map[string]T),
	}
}

// Get gets the object by the given key.
func (g *Group[T]) Get(key string) T {
	g.RLock()
	v, ok := g.vals[key]
	if ok {
		g.RUnlock()
		return v
	}
	g.RUnlock()

	// slowpath for group don`t have specified key value
	g.Lock()
	defer g.Unlock()
	v, ok = g.vals[key]
	if ok {
		return v
	}
	v = g.new(key)
	g.vals[key] = v
	return v
}

// Reset resets the new function and deletes all existing objects.
func (g *Group[T]) Reset(new CreaterFunc[T]) {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	g.Lock()
	g.new = new
	g.Unlock()
	g.Clear()
}

// Clear deletes all objects.
func (g *Group[T]) Clear() {
	g.Lock()
	g.vals = make(map[string]T)
	g.Unlock()
}
