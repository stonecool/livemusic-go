package account

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/account/state"
	"github.com/stonecool/livemusic-go/internal/account/types"
	"github.com/stonecool/livemusic-go/internal/message"
)

func CreateAccount(category string) (types.Account, error) {
	acc := &account{
		Category:     category,
		msgChan:      make(chan *message.AsyncMessage),
		done:         make(chan struct{}),
		stateHandler: state.NewStateHandler(category),
	}

	return acc, nil
}

func GetAccount(id int) (types.Account, error) {
	acc, err := repo.get(id)
	if err != nil {
		return nil, err
	}

	acc.stateHandler = state.NewStateHandler(acc.Category)

	var instance types.Account
	switch acc.Category {
	case "wechat":
		instance = &WeChatAccount{acc}
	default:
		instance = acc
	}

	return instance, nil
}

func SubmitCrawlTask(category string) error {
	result := make(chan error, 1)
	task := &message.CrawlTask{
		Category: category,
		Message:  message.NewAsyncMessage(message.AccountCmd_Crawl, result),
	}

	if err := message.DefaultQueue.PushTask(task); err != nil {
		return fmt.Errorf("failed to submit task: %w", err)
	}

	return <-result
}
