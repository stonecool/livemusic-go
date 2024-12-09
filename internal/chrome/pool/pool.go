package pool

import (
	"context"
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"sync"
	"time"

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
	mu         sync.Mutex
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
		return fmt.Errorf("instance on:%s exists", chrome.GetAddr())
	}

	p.chromes[chrome.GetAddr()] = chrome

	return nil
}

func (p *pool) Login(chrome types.Chrome, cat string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	instance, exists := p.chromes[chrome.GetAddr()]
	if !exists {
		internal.Logger.Error("instance not exists in pool",
			zap.Int("chromeID", chrome.GetID()),
			zap.String("addr", chrome.GetAddr()))
		return
	}

	category, ok := p.categories[cat]
	if ok {
		if category.ContainChrome(chrome.GetAddr()) {
			internal.Logger.Warn("instance already in category",
				zap.Int("chromeID", chrome.GetID()),
				zap.String("category", cat))
			return
		}
	}

	acc, err := account.GetAccount(chrome.GetID())
	if err != nil {
		internal.Logger.Error("failed to get account",
			zap.Error(err),
			zap.Int("chromeID", chrome.GetID()))
		return
	}

	ctx, cancel := instance.GetNewContext()
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	err = chromedp.Run(ctx,
		util.GetQRCode(acc),
		acc.WaitLogin(),
		util.SaveCookies(acc),
		chromedp.Stop(),
	)

	if err != nil {
		internal.Logger.Error("failed to run chrome commands",
			zap.Error(err),
			zap.Int("chromeID", chrome.GetID()),
			zap.String("category", cat))
		return
	}

	category.AddChrome(instance)
}

func (p *pool) GetChromesByCategory(cat string) []types.Chrome {
	p.mu.Lock()
	defer p.mu.Unlock()

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
		if err := instance.ExecuteTask(message.ITask); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no available account found for category: %s", category)
}

func (p *pool) GetAllChromes() map[string]types.Chrome {
	return p.chromes
}
