package chrome

import (
	"sync"
)

type Category struct {
	name    string
	chromes map[int]*Chrome
	mu      sync.RWMutex
}

func newCategory(name string) *Category {
	return &Category{
		name:    name,
		chromes: make(map[int]*Chrome),
	}
}

func (c *Category) AddChrome(chrome *Chrome) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.chromes[chrome.ID]; ok {
		// TODO log
		return
	}

	c.chromes[chrome.ID] = chrome
}

func (c *Category) GetChromes() []*Chrome {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*Chrome, len(c.chromes))
	for _, ins := range c.chromes {
		result = append(result, ins)
	}

	return result
}

func (c *Category) ContainChrome(id int) bool {
	_, ok := c.chromes[id]
	return ok
}
