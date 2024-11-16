package chrome

import (
	"fmt"

	"gorm.io/gorm"
)

type Factory interface {
	CreateChrome(ip string, port int, debuggerURL string) (*Chrome, error)
	GetChrome(id int) (*Chrome, error)
}

type factoryImpl struct {
	repo Repository
}

func NewFactory(repo Repository) Factory {
	return &factoryImpl{repo: repo}
}

func (f *factoryImpl) CreateChrome(ip string, port int, debuggerURL string) (*Chrome, error) {
	// 创建 model 实例
	m := &model{
		IP:          ip,
		Port:        port,
		DebuggerURL: debuggerURL,
		State:       int(STATE_UNINITIALIZED),
	}

	// 验证实例
	v := NewValidator()
	if err := v.ValidateChrome(&Chrome{
		IP:          m.IP,
		Port:        m.Port,
		DebuggerURL: m.DebuggerURL,
	}); err != nil {
		return nil, fmt.Errorf("invalid chrome instance: %w", err)
	}

	chrome := m.toEntity()
	if err := f.repo.Create(chrome); err != nil {
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	if err := chrome.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize chrome instance: %w", err)
	}

	return chrome, nil
}

func (f *factoryImpl) GetChrome(id int) (*Chrome, error) {
	chrome, err := f.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get chrome instance: %w", err)
	}

	// 如果实例需要重新初始化，则进行初始化
	if chrome.NeedsReInitialize() {
		if err := chrome.RetryInitialize(3); err != nil {
			return nil, fmt.Errorf("failed to reinitialize chrome instance: %w", err)
		}
	}

	return chrome, nil
}

// 便捷创建方法
func CreateInstance1(db *gorm.DB, ip string, port int, debuggerURL string) (*Chrome, error) {
	repo := NewRepositoryDB(db)
	factory := NewFactory(repo)
	return factory.CreateChrome(ip, port, debuggerURL)
}

//// 新增便捷获取方法
//func GetInstance(db *gorm.DB, id int) (*Chrome, error) {
//	repo := NewRepositoryDB(db)
//	factory := NewFactory(repo)
//	return factory.GetChrome(id)
//}
