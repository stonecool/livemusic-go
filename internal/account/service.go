package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/database"
)

var accountCache *cache.Memo
var repo database.Repository[*model]

func init() {
	accountCache = cache.New(func(id int) (interface{}, error) {
		return getAccount(id)
	})
	repo = newRepositoryDB(database.DB)
}

func createAccount(account *Account) IAccount {
	switch account.Category {
	case "wechat":
		wechatAccount := &WeChatAccount{Account: account}
		return wechatAccount
	default:
		return account
	}
}


func CreateAccount(category string) (IAccount, error) {
	m := &model{Category: category}
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := repo.Create(m); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	account := createAccount(m.toEntity())
	account.Init()

	if err := accountCache.Set(account.GetID(), account); err != nil {
		return nil, fmt.Errorf("failed to cache account: %w", err)
	}

	return account, nil
}

func getAccount(id int) (IAccount, error) {
	m, err := repo.Get(id)
	if err != nil {
		return nil, err
	}

	account := createAccount(m.toEntity())
	account.Init()
	return account, nil
}


func GetAccount(ID int) (IAccount, error) {
	if acc, err := accountCache.Get(ID); err != nil {
		return nil, err
	} else {
		return acc.(IAccount), nil
	}
}

