package pool

import (
	"sync"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"go.uber.org/zap"
)

type category struct {
	name    string
	chromes sync.Map
}

func newCategory(name string) *category {
	return &category{
		name: name,
	}
}

func (c *category) AddChrome(chrome types.Chrome) {
	if _, loaded := c.chromes.LoadOrStore(chrome.GetAddr(), chrome); loaded {
		internal.Logger.Warn("chrome already exists in category",
			zap.String("category", c.name),
			zap.Int("chromeID", chrome.GetID()),
			zap.String("addr", chrome.GetAddr()))
	}
}

func (c *category) GetChromes() []types.Chrome {
	var result []types.Chrome
	c.chromes.Range(func(key, value interface{}) bool {
		result = append(result, value.(types.Chrome))
		return true
	})
	return result
}

func (c *category) ContainChrome(addr string) bool {
	_, ok := c.chromes.Load(addr)
	return ok
}
