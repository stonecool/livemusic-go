package chrome

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account"
	"go.uber.org/zap"
)

func SetCookies(account account.IAccount) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(account.GetCookies()); err != nil {
			internal.Logger.Error("failed to unmarshal cookies",
				zap.Error(err),
				zap.Int("accountID", account.GetID()))
			return
		}

		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

func SaveCookies(account account.IAccount) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			internal.Logger.Error("failed to get cookies",
				zap.Error(err),
				zap.Int("accountID", account.GetID()))
			return
		}

		data, err := network.GetCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			internal.Logger.Error("failed to marshal cookies",
				zap.Error(err),
				zap.Int("accountID", account.GetID()))
			return
		}

		var url string
		err = chromedp.Evaluate(`window.location.href`, &url).Do(ctx)
		if err != nil {
			internal.Logger.Error("failed to get current URL",
				zap.Error(err),
				zap.Int("accountID", account.GetID()))
			return
		}

		account.SetLastURL(url)
		if err := account.SaveCookies(data); err != nil {
			internal.Logger.Error("failed to save cookies",
				zap.Error(err),
				zap.Int("accountID", account.GetID()))
		}

		return
	}
}
