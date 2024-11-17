package chrome

import (
	"sync"
)

type category struct {
	name    string
	chromes map[int]*Chrome
	mu      sync.RWMutex
}

func newCategory(name string) *category {
	return &category{
		name:    name,
		chromes: make(map[int]*Chrome),
	}
}

func (c *category) AddChrome(chrome *Chrome) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.chromes[chrome.ID]; ok {
		// TODO log
		return
	}

	c.chromes[chrome.ID] = chrome
}

func (c *category) GetChromes() []*Chrome {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*Chrome, len(c.chromes))
	for _, ins := range c.chromes {
		result = append(result, ins)
	}

	return result
}

func (c *category) ContainChrome(id int) bool {
	_, ok := c.chromes[id]
	return ok
}
