package factory

import (
	"github.com/stonecool/livemusic-go/internal/crawlaccount"
)

// ICrawlFactory 抽象工厂接口
type ICrawlFactory interface {
	CreateAccount() crawlaccount.ICrawlAccount
}
