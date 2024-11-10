package chrome

import (
	"fmt"
	"sync"
)

// 全局唯一的实例池
var globalPool *InstancePool

type InstancePool struct {
	instances  map[string]*Instance
	categories map[string]*Category
	mu         sync.Mutex
}

// 在包初始化时创建实例池
func init() {
	globalPool = &InstancePool{
		instances:  make(map[string]*Instance),
		categories: make(map[string]*Category),
	}
}

// 获取全局实例池
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

	if _, exists := ip.instances[ins.getAddr()]; exists {
		fmt.Printf("instance on:%s exists", ins.getAddr())
		return nil, nil
	}

	ip.instances[ins.getAddr()] = ins
	for cat := range ins.getAccounts() {
		if _, exists := ip.categories[cat]; !exists {
			ip.categories[cat] = &Category{
				name: cat,
			}
		}
		ip.categories[cat].AddInstance(ins)
	}

	return ins, nil
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

func (ip *InstancePool) ExecuteTask(cat string, task func(instance *Instance) error) error {
	for _, ins := range ip.GetInstancesByCategory(cat) {
		if ins.isAvailable(cat) {
			err := task(ins)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("no instance available for category: %s", cat)
}
