package main

import "sync"

const MAX_PRIORITY float64 = 3.0
const MIN_PRIORITY float64 = 0.0

// Holds and managers priorities for node connections
type PriorityManager struct {
	lock  sync.Mutex
	inner map[NodeCon]float64
}

func NewPriorityManager() *PriorityManager {
	m := new(PriorityManager)
	m.inner = make(map[NodeCon]float64)
	return m
}

func (self *PriorityManager) GetPriority(from Node, to Node) float64 {
	self.lock.Lock()
	defer self.lock.Unlock()

	con := NodeCon{from, to}
	value, ok := self.inner[con]
	if !ok {
		value = 1.0
		self.inner[con] = value
	}

	return value
}

// Updates priority for node connection
// NOTE: Unsafe to use concurrently
func (self *PriorityManager) updatePriorityUnsafe(
	from Node,
	to Node,
	delta float64,
) {
	// Since we're dealing with undirected graph, we need to update
	// connections for both directions
	connections := [2]NodeCon{{from, to}, {to, from}}
	for _, con := range connections {
		value, ok := self.inner[con]
		if !ok {
			value = 1.0
		}
		value += delta
		// Priority 0.0 would result it weight 0, which would result in this
		// connection never being selected at all. Not what we want.
		// Priority less than 0.0 would result in a crash.
		if value <= MIN_PRIORITY {
			value = 0.000001
		} else if value > MAX_PRIORITY {
			value = MAX_PRIORITY
		}
		self.inner[con] = value
	}
}

// Updates priority for the entire path
func (self *PriorityManager) UpdatePriorityForPath(path Path, delta float64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if len(path) < 2 {
		panic("Invalid path, len is less than 2")
	}

	for i := 0; i < len(path)-1; i++ {
		self.updatePriorityUnsafe(path[i], path[i+1], delta)
	}
}
