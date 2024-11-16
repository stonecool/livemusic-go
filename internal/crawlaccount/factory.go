package crawlaccount

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/config"
)

type factoryImpl struct {
	repo IRepository
}

func NewFactory(repo IRepository) IFactory {
	return &factoryImpl{repo: repo}
}

func (f *factoryImpl) CreateAccount(category string) (*Account, error) {
	// 检查账号类型是否支持
	if _, ok := config.AccountMap[category]; !ok {
		return nil, fmt.Errorf("unsupported crawlaccount category: %s", category)
	}

	// 创建账号模型
	model := &accountModel{
		Category: category,
	}

	account := model.toEntity()
	if err := f.repo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create crawlaccount: %w", err)
	}

	account.Init()
	return account, nil
}
