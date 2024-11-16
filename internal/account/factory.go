package account

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/config"
)

type Factory interface {
	CreateAccount(category string) (Account, error)
}

type factoryImpl struct {
	repo Repository
}

func NewFactory(repo Repository) *Factory {
	return &factoryImpl{repo: repo}
}

func (f *Factory) CreateAccount(category string) (Account, error) {
	// 检查账号类型是否支持
	if _, ok := config.AccountMap[category]; !ok {
		return nil, fmt.Errorf("unsupported account category: %s", category)
	}

	// 创建账号模型
	model := &accountModel{
		Category: category,
		State:    AS_EXPIRED, // 初始状态为过期
	}

	// 创建账号实例
	account := NewAccount(model)

	// 持久化到数据库
	if err := f.repo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}
