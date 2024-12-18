package account

import (
	"github.com/stonecool/livemusic-go/internal/account/state"
	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stonecool/livemusic-go/internal/message"
)

var (
	repo repository
)

func init() {
	repo = newRepositoryDB(database.DB)
}

func CreateAccount(category string) (IAccount, error) {
	acc := &account{
		Category:     category,
		msgChan:      make(chan *message.AsyncMessage),
		done:         make(chan struct{}),
		stateHandler: state.NewStateHandler(category),
	}

	return acc, nil
}

func GetAccount(id int) (IAccount, error) {
	acc, err := repo.get(id)
	if err != nil {
		return nil, err
	}

	acc.stateManager = state.selectStateManager(acc.Category)

	var instance IAccount
	switch acc.Category {
	case "wechat":
		instance = &WeChatAccount{acc}
	default:
		instance = acc
	}

	return instance, nil
}
