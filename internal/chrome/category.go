package chrome

import (
	"sync"
)

type Category struct {
	name      string
	instances map[int]*Instance
	mu        sync.RWMutex
}

func newCategory(name string) *Category {
	return &Category{
		name:      name,
		instances: make(map[int]*Instance),
	}
}

func (c *Category) AddInstance(ins *Instance) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.instances[ins.ID]; ok {
		// TODO log
		return
	}

	c.instances[ins.ID] = ins
}

func (c *Category) GetInstances() []*Instance {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*Instance, len(c.instances))
	for _, ins := range c.instances {
		result = append(result, ins)
	}

	return result
}

func (c *Category) ContainInstance(id int) bool {
	_, ok := c.instances[id]
	return ok
}
