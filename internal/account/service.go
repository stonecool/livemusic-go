package account

import (
	"github.com/stonecool/livemusic-go/internal/database"
)

var (
	repo repository
)

func init() {
	repo = newRepositoryDB(database.DB)
}

func CreateAccount(category string) (IAccount, error) {
	acc, err := repo.create(category, stateNew)
	if err != nil {
		return nil, err
	}

	return GetAccount(acc.ID)
}

func GetAccount(id int) (IAccount, error) {
	acc, err := repo.get(id)
	if err != nil {
		return nil, err
	}

	acc.stateManager = selectStateManager(acc.Category)

	var instance IAccount
	switch acc.Category {
	case "wechat":
		instance = &WeChatAccount{acc}
	default:
		instance = acc
	}

	return instance, nil
}
