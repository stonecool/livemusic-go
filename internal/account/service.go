package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stonecool/livemusic-go/internal/message"
	"time"
)

var (
	accountCache *cache.Memo
	repo         repository
)

func init() {
	accountCache = cache.New(func(id int) (interface{}, error) {
		return getAccount(id)
	})
	repo = newRepositoryDB(database.DB)
}

func CreateAccount(category string) (IAccount, error) {
	acc, err := repo.create(category, stateNew)
	if err != nil {
		return nil, err
	}

	return getAccount(acc.ID)
}

func getAccount(id int) (IAccount, error) {
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

	go func() {
		select {
		case instance.GetMsgChan() <- message.NewAsyncMessageWithCmd(message.CrawlCmd_Initial, nil):
		case <-time.After(5 * time.Second):
			fmt.Println("time out")
		}
	}()

	return instance, nil
}

func GetAccount(id int) (IAccount, error) {
	if instance, err := accountCache.Get(id); err != nil {
		return nil, err
	} else {
		return instance.(IAccount), nil
	}
}
