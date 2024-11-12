package crawl

import (
	"context"
	"github.com/stonecool/livemusic-go/internal/config"
)

// 微信爬虫
type WeChatCrawl struct {
	state   CrawlState
	config  *config.Account
	ch      chan *ClientMessage
	context context.Context
}
