package pool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stonecool/livemusic-go/internal"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/chrome/util"
	"github.com/stonecool/livemusic-go/internal/message"
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

func (p *pool) Login(chrome types.Chrome, cat string) {
	category, acc, err := p.prepareLogin(chrome, cat)
	if err != nil {
		internal.Logger.Error("failed to prepare login",
			zap.Error(err),
			zap.Int("chromeID", chrome.GetID()),
			zap.String("category", cat))
		return
	}

	if err := p.doLogin(chrome, acc); err != nil {
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

func (p *pool) prepareLogin(chrome types.Chrome, cat string) (*category, account.IAccount, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, exists := p.chromes[chrome.GetAddr()]
	if !exists {
		return nil, nil, fmt.Errorf("instance not exists in pool: %s", chrome.GetAddr())
	}

	category, ok := p.categories[cat]
	if ok && category.ContainChrome(chrome.GetAddr()) {
		return nil, nil, fmt.Errorf("instance already in category: %s", cat)
	}

	acc, err := account.GetAccount(chrome.GetID())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account: %w", err)
	}

	return category, acc, nil
}

func (p *pool) doLogin(chrome types.Chrome, acc account.IAccount) error {
	ctx, cancel := chrome.GetNewContext()
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	return chromedp.Run(ctx,
		util.GetQRCode(acc),
		acc.WaitLogin(),
		util.SaveCookies(acc),
		chromedp.Stop(),
	)
}

func (p *pool) GetChromesByCategory(cat string) []types.Chrome {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if cat, exists := p.categories[cat]; exists {
		return cat.GetChromes()
	} else {
		return nil
	}
}

func (p *pool) DispatchTask(category string, message *message.AsyncMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 获取该分类下的所有实例
	instances := p.GetChromesByCategory(category)
	if len(instances) == 0 {
		return fmt.Errorf("no instance available for category: %s", category)
	}

	// 遍历实例找到可用的账号
	for _, instance := range instances {
		if !instance.IsAvailable() {
			continue
		}
		
		if err := instance.ExecuteTask(message.ITask); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no available account found for category: %s", category)
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
