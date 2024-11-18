package account

import (
	"fmt"
)

type factory struct {
	repo IRepository
}

func newAccountFactory(repo IRepository) *factory {
	return &factory{repo: repo}
}

func (f *factory) Create(category string) (IAccount, error) {
	v := NewValidator()
	if err := v.validateCategory(category); err != nil {
		return nil, fmt.Errorf("invalid account category: %w", err)
	}

	model := &model{Category: category}
	account := model.toEntity()

	if err := v.ValidateAccount(account); err != nil {
		return nil, fmt.Errorf("invalid account: %w", err)
	}

	if err := f.repo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return f.createInstance(account)
}

func (f *factory) Get(id int) (IAccount, error) {
	account, err := f.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return f.createInstance(account)
}

func (f *factory) createInstance(account *Account) (IAccount, error) {
	switch account.Category {
	case "wechat":
		wechatAccount := &WeChatAccount{Account: account}
		wechatAccount.Init()
		return wechatAccount, nil
	default:
		return account, nil
	}
}
