package chrome

import (
	"fmt"
)


type Factory struct {
	repo repository
}

func NewFactory(repo repository) *Factory {
	return &Factory{repo: repo}
}

func (f *Factory) CreateChrome(ip string, port int, debuggerURL string) (*Chrome, error) {
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
