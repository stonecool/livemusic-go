package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/database"
)

var (
	accountCache *cache.Memo
	accountRepo  repository
)

func init() {
	accountCache = cache.New(func(id int) (interface{}, error) {
		return getAccount(id)
	})
	accountRepo = newRepositoryDB(database.DB)
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
	account, err := accountRepo.Create(category)
	if err != nil {
		return nil, err
	}

	instance := createAccount(account)
	instance.Init()

	if err := accountCache.Set(account.GetID(), account); err != nil {
		return nil, fmt.Errorf("failed to cache account: %w", err)
	}

	return instance, nil
}

func getAccount(id int) (IAccount, error) {
	account, err :=  accountRepo.Get(id)
	if err != nil {
		return nil, err
	}

	instance := createAccount(account)
	instance.Init()

	return instance, nil
}


func GetAccount(ID int) (IAccount, error) {
	if acc, err := accountCache.Get(ID); err != nil {
		return nil, err
	} else {
		return acc.(IAccount), nil
	}
}

