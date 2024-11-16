package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account"
	"log"
)

// SetCookies
func SetCookies(account account.ICrawlAccount) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(account.GetCookies()); err != nil {
			log.Printf("%s", err)
			return
		}

		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

// SaveCookies
func SaveCookies(account account.ICrawlAccount) chromedp.ActionFunc {
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

		account.SetLastURL(url)
		if err := account.SaveCookies(data); err != nil {
			fmt.Printf("%v", err)
		}

		return
	}
}
