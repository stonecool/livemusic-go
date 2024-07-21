package crawl

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/cache"
)

var msgCache *cache.Memo

func init() {
	msgCache = cache.New(getCrawlMsg)
}

func getCrawlMsg(id int) (interface{}, error) {
	msg := &internal.CrawlMsg{ID: id}

	if err := msg.Get(); err != nil {
		return nil, err
	}

	switch msg.DataType {
	case "Livehouse":
		return nil, error(nil)

	}

	// TODO
	return nil, nil
}

func GetCrawlMsg(id int) (*internal.CrawlMsg, error) {
	msg, err := msgCache.Get(id)
	if err != nil {
		return nil, err
	} else {
		return msg.(*internal.CrawlMsg), nil
	}
}
