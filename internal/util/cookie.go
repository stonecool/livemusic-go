package util

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/crawl"
	"log"
)

// getCookies
func getCookies(iCrawl crawl.ICrawl) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(iCrawl.GetCookies()); err != nil {
			log.Printf("%s", err)
			return
		}

		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

// setCookies
func setCookies(iCrawl crawl.ICrawl) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return
		}

		data, err := network.GetCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}

		iCrawl.SaveCookies(data)

		return
	}
}
