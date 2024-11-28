package pool

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"go.uber.org/zap"
	"sync"
)

type category struct {
	name    string
	chromes map[string]*instance.Chrome
	mu      sync.RWMutex
}

func newCategory(name string) *category {
	return &category{
		name:    name,
		chromes: make(map[string]*instance.Chrome),
	}
}

func (c *category) AddChrome(chrome *instance.Chrome) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.chromes[chrome.GetAddr()]; ok {
		internal.Logger.Warn("chrome already exists in category",
			zap.String("category", c.name),
			zap.Int("chromeID", chrome.ID))
		return
	}

	c.chromes[chrome.GetAddr()] = chrome
}

func (c *category) GetChromes() []*instance.Chrome {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*instance.Chrome, 0, len(c.chromes))
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
