package chrome

import (
	"sync"
)

type InstancePool struct {
	instances  map[string]*Instance
	categories map[string]*Category
	mu         sync.Mutex
}

// NewInstancePool 创建chrome实例池
func NewInstancePool() *InstancePool {
	return &InstancePool{
		instances:  make(map[string]*Instance),
		categories: make(map[string]*Category),
	}
}

// AddInstance 添加新的实例到池
func (ip *InstancePool) AddInstance(ins *Instance) {
	if _, exists := ip.instances[ins.addr]; exists {
		return
	}

	ip.mu.Lock()
	defer ip.mu.Unlock()

	ip.instances[ins.addr] = ins
	for cat := range ins.getAccounts() {
		if _, exists := ip.categories[cat]; !exists {
			ip.categories[cat] = &Category{
				name: cat,
			}
		}
		ip.categories[cat].AddInstance(ins)
	}
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

func (ip *InstancePool) ExecuteTask(cat string, task func(instance *Instance)) {
	for _, ins := range ip.GetInstancesByCategory(cat) {
		if ins.isAvailable(cat) {
			task(ins)
			break
		}
	}
}
