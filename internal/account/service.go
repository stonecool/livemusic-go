package account

import (
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/database"
)

var (
	accountCache *cache.Memo
	accountRepo  repository
)

func init() {
	accountCache = cache.New(func(id int) (interface{}, error) {
		return getInstance(id)
	})
	accountRepo = newRepositoryDB(database.DB)
}

func CreateInstance(category string) (IAccount, error) {
	account, err := accountRepo.create(category, stateNew)
	if err != nil {
		return nil, err
	}

	return getInstance(account.ID)
}

func getInstance(id int) (IAccount, error) {
	account, err := accountRepo.get(id)
	if err != nil {
		return nil, err
	}

	account.stateManager = selectStateManager(account.Category)

	var instance IAccount
	switch account.Category {
	case "wechat":
		instance = &WeChatAccount{account: account}
	default:
		instance = account
	}

	instance.Init()

	return instance, nil
}

func GetInstance(ID int) (IAccount, error) {
	if instance, err := accountCache.Get(ID); err != nil {
		return nil, err
	} else {
		return instance.(IAccount), nil
	}
}
