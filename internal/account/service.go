package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/cache"
)

var accountCache *cache.Memo

func init() {
	accountCache = cache.New(func(id int) (interface{}, error) {
		return getAccount(id)
	})
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
	account, err := CreateAccountInDB(category)
	if err != nil {
		return nil, err
	}

	instance := createAccount(account)
	instance.Init()

	if err := accountCache.Set(account.GetID(), account); err != nil {
		return nil, fmt.Errorf("failed to cache account: %w", err)
	}

	return account, nil
}

func getAccount(id int) (IAccount, error) {
	account, err := GetAccountByID(id)
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

