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
		return getAccount(id)
	})
	accountRepo = newRepositoryDB(database.DB)
}

func CreateAccount(category string) (IAccount, error) {
	acc, err := accountRepo.create(category, stateNew)
	if err != nil {
		return nil, err
	}

	return getAccount(acc.ID)
}

func getAccount(id int) (IAccount, error) {
	acc, err := accountRepo.get(id)
	if err != nil {
		return nil, err
	}

	acc.stateManager = selectStateManager(acc.Category)

	var instance IAccount
	switch acc.Category {
	case "wechat":
		instance = &WeChatAccount{account: acc}
	default:
		instance = acc
	}

	instance.Init()

	return instance, nil
}

func GetAccount(id int) (IAccount, error) {
	if acc, err := accountCache.Get(id); err != nil {
		return nil, err
	} else {
		return acc.(IAccount), nil
	}
}
