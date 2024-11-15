package chrome

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/crawlaccount"
	"github.com/stonecool/livemusic-go/internal/model"
)

// 全局唯一的实例池
var globalPool *InstancePool

type InstancePool struct {
	instances      map[int]*Instance
	addr2Instances map[string]*Instance
	categories     map[string]*Category
	mu             sync.Mutex
}

// init 在包初始化时创建实例池
func init() {
	globalPool = &InstancePool{
		instances:      make(map[int]*Instance),
		addr2Instances: make(map[string]*Instance),
		categories:     make(map[string]*Category),
	}
}

// GetPool 获取全局实例池
func GetPool() *InstancePool {
	return globalPool
}

// AddInstance 添加新的实例到池
func (ip *InstancePool) AddInstance(id int) (*Instance, error) {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	ins, err := GetInstance(id)
	if err != nil {
		return nil, err
	}

	if _, exists := ip.instances[ins.ID]; exists {
		fmt.Printf("instance on:%s exists", ins.GetAddr())
		return nil, nil
	}

	if _, exists := ip.addr2Instances[ins.GetAddr()]; exists {
		fmt.Printf("instance on:%s exists", ins.GetAddr())
		return nil, nil
	}

	ip.instances[ins.ID] = ins
	ip.addr2Instances[ins.GetAddr()] = ins
	for cat := range ins.getAccounts() {
		if _, exists := ip.categories[cat]; !exists {
			ip.categories[cat] = newCategory(cat)
		}
		ip.categories[cat].AddInstance(ins)
	}

	return ins, nil
}

func (ip *InstancePool) Login(id int, cat string) {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	instance, exists := ip.instances[id]
	if !exists {
		fmt.Printf("instance:%d not exists in pool", id)
		return
	}

	category, ok := ip.categories[cat]
	if ok {
		if category.ContainInstance(id) {
			fmt.Printf("instance:%d already in cat:%s", id, cat)
			return
		}
	}

	data := map[string]interface{}{
		"account_type": cat,
	}

	m, err := model.AddCrawlAccount(data)
	if err != nil {
		return
	}

	crawlAccount := crawlaccount.NewCrawlAccount(m)
	ctx, cancel := instance.GetNewContext()
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	//ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))

	err = chromedp.Run(ctx,
		internal.GetQRCode(crawlAccount),
		crawlAccount.WaitLogin(),
		internal.SaveCookies(crawlAccount),
		chromedp.Stop(),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	category.AddInstance(instance)

	return
}

func (ip *InstancePool) GetInstancesByCategory(cat string) []*Instance {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	if cat, exists := ip.categories[cat]; exists {
		return cat.GetInstances()
	} else {
		return nil
	}
}

func (ip *InstancePool) DispatchTask(category string, task *internal.CrawlTask) error {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	// 获取该分类下的所有实例
	instances := ip.GetInstancesByCategory(category)
	if len(instances) == 0 {
		return fmt.Errorf("no instance available for category: %s", category)
	}

	// 遍历实例找到可用的账号
	for _, instance := range instances {
		if err := instance.ExecuteTask(task); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no available account found for category: %s", category)
}
