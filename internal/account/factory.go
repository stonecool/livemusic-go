package account

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/database"

	"github.com/stonecool/livemusic-go/internal/config"
)

type factory struct {
	repo IRepository
}

func newFactory(repo IRepository) *factory {
	return &factory{repo: repo}
}

func (f *factory) createAccount(category string) (IAccount, error) {
	v := NewValidator()
	if err := v.validateCategory(category); err != nil {
		return nil, fmt.Errorf("invalid account category: %w", err)
	}

	if _, ok := config.AccountMap[category]; !ok {
		return nil, fmt.Errorf("unsupported account category: %s", category)
	}

	model := &model{Category: category}
	account := model.toEntity()

	if err := v.ValidateAccount(account); err != nil {
		return nil, fmt.Errorf("invalid account: %w", err)
	}

	if err := f.repo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	account.Init()
	switch account.Category {
	case "wechat":
		return &WeChatAccount{Account: *account}, nil
	default:
		return account, nil
	}
}

func CreateAccount(category string) (IAccount, error) {
	repo := NewRepositoryDB(database.DB)
	factory := newFactory(repo)
	return factory.createAccount(category)
}


