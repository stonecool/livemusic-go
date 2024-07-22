package crawl

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/config"
	"log"
)

var crawlCache *cache.Memo

func init() {
	crawlCache = cache.New(getCrawl)
}

func getCrawl(id int) (interface{}, error) {
	account := &internal.CrawlAccount{ID: id}
	err := account.Get()
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	cfg, ok := config.AccountMap[account.AccountType]
	if !ok {
		return nil, error(nil)
	}

	var crawl ICrawl
	switch account.AccountType {
	case "WeChat":
		crawl = &WeChatCrawl{
			Crawl: Crawl{
				Account: account,
				state:   internal.CrawlState_Uninitialized,
				config:  &cfg,
				ch:      make(chan *internal.Message),
			},
		}
	}

	go crawl.Start()
	return crawl, nil
}

func GetCrawl(id int) (*Crawl, error) {
	crawl, err := crawlCache.Get(id)
	if err != nil {
		return nil, err
	} else {
		return crawl.(*Crawl), nil
	}
}
