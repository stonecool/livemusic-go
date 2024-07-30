package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
)

// setCookies
func setCookies(iCrawl ICrawl) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(iCrawl.GetCookies()); err != nil {
			log.Printf("%s", err)
			return
		}

		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

// saveCookies
func saveCookies(iCrawl ICrawl) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		data, err := network.GetCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			fmt.Printf("%v", err)

			return
		}

		var url string
		err = chromedp.Evaluate(`window.location.href`, &url).Do(ctx)
		if err != nil {
			return
		}

		iCrawl.SetLastLoginURL(url)
		if err := iCrawl.SaveCookies(data); err != nil {
			fmt.Printf("%v", err)
		}

		return
	}
}
