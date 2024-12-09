package pool

import (
	"sync"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"go.uber.org/zap"
)

type category struct {
	name    string
	chromes map[string]types.Chrome
	mu      sync.RWMutex
}

func newCategory(name string) *category {
	return &category{
		name:    name,
		chromes: make(map[string]types.Chrome),
	}
}

func (c *category) AddChrome(chrome types.Chrome) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.chromes[chrome.GetAddr()]; ok {
		internal.Logger.Warn("chrome already exists in category",
			zap.String("category", c.name),
			zap.Int("chromeID", chrome.GetID()),
			zap.String("addr", chrome.GetAddr()))
		return
	}

	c.chromes[chrome.GetAddr()] = chrome
}

func (c *category) GetChromes() []types.Chrome {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]types.Chrome, 0, len(c.chromes))
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
