package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
)

type WeChatCrawl struct {
	Crawl
}

func (crawl *WeChatCrawl) WaitLogin() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		return chromedp.WaitVisible(`#app > div.main_bd_new`, chromedp.ByID).Do(ctx)
	}
}

func (crawl *WeChatCrawl) CheckLogin() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		var msg string
		if err = chromedp.Evaluate(`document.querySelector('#body > div > div > div > div > div > div').className`, &msg).Do(ctx); err != nil {
			return
		}

		if msg == "page_error_msg" {
			return fmt.Errorf("login error")
		}

		return
	}
}

func (crawl *WeChatCrawl) GetQRCodeSelector() string {
	return "#header > div.banner > div > div > div.login__type__container.login__type__container__scan > img"
}

func (crawl *WeChatCrawl) SetLastLoginURL(url string) {
	crawl.Account.lastLoginURL = url
}

func (crawl *WeChatCrawl) GetLastLoginURL() string {
	return crawl.Account.lastLoginURL
}
