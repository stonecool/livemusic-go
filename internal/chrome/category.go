package chrome

import (
	"sync"
)

type Category struct {
	name      string
	instances []*Instance
	mu        sync.RWMutex
}

func (c *Category) AddInstance(ins *Instance) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances = append(c.instances, ins)
}

func (c *Category) GetInstances() []*Instance {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]*Instance, len(c.instances))
	copy(result, c.instances)
	return result
}
