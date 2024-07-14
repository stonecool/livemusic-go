package crawl

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/cache"
)

var crawlMsgFactory *cache.Memo

func init() {
	crawlMsgFactory = cache.New(getCrawlMsgProducer)
}

func getCrawlMsgProducer(id int) (interface{}, error) {
	_, err := internal.GetCrawlMsg(id)
	if err != nil {
		return nil, err
	}

	// TODO
	return nil, nil
}

func GetCrawlMsgProducer(id int) (*internal.CrawlMsg, error) {
	producer, err := crawlMsgFactory.Get(id)
	if err != nil {
		return nil, err
	} else {
		return producer.(*internal.CrawlMsg), nil
	}
}
