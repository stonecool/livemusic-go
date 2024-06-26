package crawl

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

type WxCrawl struct {
	Crawl
}

func (c *WxCrawl) Login() (bool, error) {
	//var lastLoggedPath = "https://mp.weixin.qq.com/cgi-bin/home?t=home/index&lang=zh_CN&token=2098303583"
	//get cookie to colly

	err := QRCodeLogin(c) //"#header > div.banner > div > div > div.login__type__container.login__type__container__scan > img",

	if err != nil {
		return false, err
	}

	return false, nil
}

func (c *WxCrawl) WaitLogin() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if err = chromedp.WaitVisible(`#app > div.main_bd_new`, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		return
	}
}

func (c *WxCrawl) CheckLogin() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		chromedp.Navigate(c.config.LoginURL)

		var msg string
		if err = chromedp.Evaluate(`document.querySelector('#body > div > div > div > div > div > div').className`, &msg).Do(ctx); err != nil {
			return
		}

		if msg == "page_error_msg" {
			log.Println("no login.")
			return

		}

		if err = chromedp.Evaluate(`window.location.href`, c.config.LoginURL).Do(ctx); err != nil {
			return
		}

		log.Println("login succeed.")

		return
	}
}
