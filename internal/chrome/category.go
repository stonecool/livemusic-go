package chrome

import (
	"sync"

	"github.com/stonecool/livemusic-go/internal"
	"go.uber.org/zap"
)

type category struct {
	name    string
	chromes map[string]*Chrome
	mu      sync.RWMutex
}

func newCategory(name string) *category {
	return &category{
		name:    name,
		chromes: make(map[string]*Chrome),
	}
}

func (c *category) AddChrome(chrome *Chrome) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.chromes[chrome.getAddr()]; ok {
		internal.Logger.Warn("chrome already exists in category",
			zap.String("category", c.name),
			zap.Int("chromeID", chrome.ID))
		return
	}

	c.chromes[chrome.getAddr()] = chrome
}

func (c *category) GetChromes() []*Chrome {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*Chrome, 0, len(c.chromes))
	for _, chrome := range c.chromes {
		result = append(result, chrome)
	}

	return result
}

func (c *category) ContainChrome(addr string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.chromes[addr]
	return ok
}
