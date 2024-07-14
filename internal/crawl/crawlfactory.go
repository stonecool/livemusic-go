package crawl

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/cache"
	"log"
	"reflect"
)

var crawlFactory *cache.Memo

func init() {
	crawlFactory = cache.New(getCrawl)
}

// getCrawl
func getCrawl(id int) (interface{}, error) {
	account, err := internal.GetCrawlAccount(id)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	// FIXME Is this check necessary?
	if reflect.ValueOf(*account).IsZero() {
		return &Crawl{}, nil
	}

	var crawl ICrawl
	switch account.AccountType {
	case "wx":
		crawl = &WxCrawl{
			Crawl: Crawl{
				Account: account,
				//config:  &cfg,
				ch: make(chan *internal.Message),
			},
		}
	}

	go crawl.Start()
	return crawl, nil
}

// GetCrawl
func GetCrawl(id int) (*Crawl, error) {
	crawl, err := crawlFactory.Get(id)
	if err != nil {
		return nil, err
	} else {
		return crawl.(*Crawl), nil
	}
}
