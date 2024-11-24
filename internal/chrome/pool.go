package chrome

import (
	"context"
	"fmt"
	"github.com/stonecool/livemusic-go/internal/message"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account"
)

// 全局唯一的实例池
var globalPool *Pool

type Pool struct {
	chromes      map[int]*Chrome
	addr2Chromes map[string]*Chrome
	categories   map[string]*category
	mu           sync.Mutex
}

// init 在包初始化时创建实例池
func init() {
	globalPool = &Pool{
		chromes:      make(map[int]*Chrome),
		addr2Chromes: make(map[string]*Chrome),
		categories:   make(map[string]*category),
	}
}

// GetPool 获取全局实例池
func GetPool() *Pool {
	return globalPool
}

// AddChrome 添加新的实例到池
func (p *Pool) AddChrome(id int) (*Chrome, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ins, err := GetInstance(id)
	if err != nil {
		return nil, err
	}

	if _, exists := p.chromes[ins.ID]; exists {
		fmt.Printf("instance on:%s exists", ins.GetAddr())
		return nil, nil
	}

	if _, exists := p.addr2Chromes[ins.GetAddr()]; exists {
		fmt.Printf("instance on:%s exists", ins.GetAddr())
		return nil, nil
	}

	p.chromes[ins.ID] = ins
	p.addr2Chromes[ins.GetAddr()] = ins
	for cat := range ins.getAccounts() {
		if _, exists := p.categories[cat]; !exists {
			p.categories[cat] = newCategory(cat)
		}
		p.categories[cat].AddChrome(ins)
	}

	return ins, nil
}

func (p *Pool) Login(id int, cat string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	instance, exists := p.chromes[id]
	if !exists {
		fmt.Printf("instance:%d not exists in pool", id)
		return
	}

	category, ok := p.categories[cat]
	if ok {
		if category.ContainChrome(id) {
			fmt.Printf("instance:%d already in cat:%s", id, cat)
			return
		}
	}

	acc, err := account.GetAccount(id)
	if err != nil {
		return
	}

	ctx, cancel := instance.GetNewContext()
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	//ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))

	err = chromedp.Run(ctx,
		GetQRCode(acc),
		acc.WaitLogin(),
		SaveCookies(acc),
		chromedp.Stop(),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	category.AddChrome(instance)

	return
}

func (p *Pool) GetChromesByCategory(cat string) []*Chrome {
	p.mu.Lock()
	defer p.mu.Unlock()

	if cat, exists := p.categories[cat]; exists {
		return cat.GetChromes()
	} else {
		return nil
	}
}

func (p *Pool) DispatchTask(category string, message *message.AsyncMessage) error {
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
