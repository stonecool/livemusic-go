package internal

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
)

// getCookies
func getCookies(iCrawl ICrawl) chromedp.ActionFunc {
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
func setCookies(iCrawl ICrawl) chromedp.ActionFunc {
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
