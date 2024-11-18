package account

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/config"
)

type factory struct {
	repo IRepository
}

func newFactory(repo IRepository) *factory {
	return &factory{repo: repo}
}

func (f *factory) createAccount(category string) (*Account, error) {
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

	return account, nil
}
