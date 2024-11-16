package crawlaccount

import (
	"github.com/stonecool/livemusic-go/internal/config"
)

type WeChatFactory struct{}

func (f *WeChatFactory) CreateAccount(cfg *config.AccountConfig) ICrawlAccount {
	return &WeChatAccount{
		CrawlAccount: CrawlAccount{
			Category: "wechat",
			State:    internal.AccountState_Uninitialized,
			msgChan:  make(chan *internal.AsyncMessage),
			done:     make(chan struct{}),
		},
		config: cfg,
	}
}
