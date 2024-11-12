package factory

import (
	"github.com/stonecool/livemusic-go/internal/config"
)

type WeChatFactory struct{}

func (f *WeChatFactory) CreateCrawl(config *config.Account) ICrawl {
	return &WeChatCrawl{
		state:  CrawlState_Uninitialized,
		config: config,
		ch:     make(chan *ClientMessage),
	}
}

func (f *WeChatFactory) CreateAccount() IAccount {
	return &WeChatAccount{}
}
