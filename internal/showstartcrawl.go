package internal

import "github.com/chromedp/chromedp"

type ShowStartCrawl struct {
	Crawl
}

func (crawl *ShowStartCrawl) GoCrawl(callback Callback) chromedp.ActionFunc {
	return nil
}

func (crawl *ShowStartCrawl) callback(ret map[string]interface{}) {

}
