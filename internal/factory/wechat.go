package factory

import (
	"github.com/stonecool/livemusic-go/internal/crawlaccount"
)

type WeChatFactory struct{}

func (f *WeChatFactory) CreateAccount() crawlaccount.ICrawlAccount {
	return &crawlaccount.WeChatAccount{}
}
