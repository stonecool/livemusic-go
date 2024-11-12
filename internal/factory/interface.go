package factory

import "github.com/stonecool/livemusic-go/internal/config"

// ICrawlFactory 抽象工厂接口
type ICrawlFactory interface {
	CreateCrawl(config *config.Account) ICrawl
	CreateAccount() IAccount
}
