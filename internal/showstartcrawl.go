package internal

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/crawl"
)

type ShowStartCrawl struct {
	crawl.Crawl
}

func (crawl *ShowStartCrawl) GoCrawl(callback Callback) chromedp.ActionFunc {
	return nil
}

func (crawl *ShowStartCrawl) callback(ret map[string]interface{}) {

}
