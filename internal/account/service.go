package account

import (
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stonecool/livemusic-go/internal/message"
	"time"
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
		instance = &WeChatAccount{acc}
	default:
		instance = acc
	}

	msg := message.NewAsyncMessage(
		&message.Message{
			Cmd: message.CrawlCmd_Initial,
		}, nil)

	go func() {
		select {
		case instance.GetMsgChan() <- msg:
			// 消息发送成功
		case <-time.After(5 * time.Second):
			// 处理发送超时
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
