package chrome

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/instance"
	"sync"
)

type InstancePool struct {
	instances  map[string]*instance.Instance
	categories map[string]*Category
	mu         sync.Mutex
}

var Pool = newInstancePool()

// newInstancePool 创建chrome实例池
func newInstancePool() *InstancePool {
	return &InstancePool{
		instances:  make(map[string]*instance.Instance),
		categories: make(map[string]*Category),
	}
}

// AddInstance 添加新的实例到池
func (ip *InstancePool) AddInstance(ins *instance.Instance) {
	if _, exists := ip.instances[ins.getAddr()]; exists {
		fmt.Printf("instance on:%s exists", ins.getAddr())
		return
	}

	ip.mu.Lock()
	defer ip.mu.Unlock()

	ip.instances[ins.getAddr()] = ins
	for cat := range ins.getAccounts() {
		if _, exists := ip.categories[cat]; !exists {
			ip.categories[cat] = &Category{
				name: cat,
			}
		}
		ip.categories[cat].AddInstance(ins)
	}
}

func (ip *InstancePool) GetInstancesByCategory(cat string) []*instance.Instance {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	if cat, exists := ip.categories[cat]; exists {
		return cat.GetInstances()
	} else {
		return nil
	}
}

func (ip *InstancePool) ExecuteTask(cat string, task func(instance *instance.Instance)) {
	for _, ins := range ip.GetInstancesByCategory(cat) {
		if ins.isAvailable(cat) {
			task(ins)
			break
		}
	}
}
