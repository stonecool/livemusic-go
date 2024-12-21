package pool

import (
	"fmt"
	"sync"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"go.uber.org/zap"
)

var GlobalPool *pool

type pool struct {
	chromes    map[string]types.Chrome
	categories map[string]*category
	mu         sync.RWMutex
}

// init 在包初始化时创建实例池
func init() {
	GlobalPool = &pool{
		chromes:    make(map[string]types.Chrome),
		categories: make(map[string]*category),
	}
}

// GetPool 获取全局实例池
func GetPool() *pool {
	return GlobalPool
}

func (p *pool) AddChrome(chrome types.Chrome) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.chromes[chrome.GetAddr()]; exists {
		err := fmt.Errorf("instance%s exists in pool", chrome.GetAddr())
		internal.Logger.Error(err.Error())
		return err
	}

	p.chromes[chrome.GetAddr()] = chrome
	return nil
}

func (p *pool) GetAllChromes() map[string]types.Chrome {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.chromes
}

func (p *pool) GetChrome(addr string) types.Chrome {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.chromes[addr]
}

func (p *pool) Login(chrome types.Chrome, cat string) {
	p.mu.RLock()

	_, exists := p.chromes[chrome.GetAddr()]
	if !exists {
		fmt.Errorf("instance not exists in pool: %s", chrome.GetAddr())
		return
	}

	category, ok := p.categories[cat]
	if ok && category.ContainChrome(chrome.GetAddr()) {
		fmt.Errorf("instance already in category: %s", cat)
		return
	}

	acc, err := account.GetAccount(chrome.GetID())
	if err != nil {
		fmt.Errorf("failed to get account: %w", err)
		return
	}
	p.mu.RUnlock()

	if err := chrome.Login(acc); err != nil {
		internal.Logger.Error("failed to login",
			zap.Error(err),
			zap.Int("chromeID", chrome.GetID()),
			zap.String("category", cat))
		return
	}

	p.mu.Lock()
	category.AddChrome(chrome)
	p.mu.Unlock()
}
