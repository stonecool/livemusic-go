package util

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func WxQRCodeLogin(loginURL, cookiesFilePath, saveCodePath string, homeURL *string) error {
	err := QRCodeLogin(loginURL, cookiesFilePath, saveCodePath,
		"#header > div.banner > div > div > div.login__type__container.login__type__container__scan > img",
		`#app > div.main_bd_new`,
		wxCheckLogin(homeURL),
	)
	if err != nil {
		return err
	}

	return nil
}

// wxCheckLogin
func wxCheckLogin(homeURL *string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		chromedp.Navigate(*homeURL)

		var url string
		if err = chromedp.Evaluate(`document.querySelector('#body > div > div > div > div > div > div').className`, &url).Do(ctx); err != nil {
			return
		}

		if url == "page_error_msg" {
			log.Println("no login.")
			return

		}

		if err = chromedp.Evaluate(`window.location.href`, homeURL).Do(ctx); err != nil {
			return
		}

		log.Println("login succeed.")
		chromedp.Stop()

		return
	}
}
