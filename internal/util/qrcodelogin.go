package util

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"time"
)

// QRCodeLogin
func QRCodeLogin(loginURL string, cookieFilePath, qrCodeFilePath string, qrcodeSel, loggedSel string, checkLoginStatus chromedp.ActionFunc) error {
	ctx, _ := chromedp.NewExecAllocator(
		context.Background(),

		append(
			chromedp.DefaultExecAllocatorOptions[:],
			//chromedp.NoDefaultBrowserCheck,
			chromedp.Flag("headless", false),
			//chromedp.Flag("hide-scrollbars", false),
			//chromedp.Flag("mute-audio", false),
			//chromedp.Flag("ignore-certificate-errors", true),
			//chromedp.Flag("disable-web-security", true),
			//chromedp.Flag("disable-gpu", false),
			//chromedp.NoFirstRun,
			//chromedp.Flag("enable-automation", false),
			//chromedp.Flag("disable-extensions", false),
		)...,
	)

	// create chrome instance
	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		loadCookies(cookieFilePath),
		checkLoginStatus,
		saveCode(loginURL, qrcodeSel, qrCodeFilePath),
		saveCookies(cookieFilePath, loggedSel),
		checkLoginStatus,
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func loadCookies(path string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if path == "" {
			return
		}

		if _, _err := os.Stat(path); os.IsNotExist(_err) {
			return
		}

		cookiesData, err := os.ReadFile(path)
		if err != nil {
			log.Printf("%s", err)
			return
		}

		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			log.Printf("%s", err)
			return
		}

		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

func saveCookies(path, sel string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if path == "" || sel == "" {
			return
		}

		if err = chromedp.WaitVisible(sel, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return
		}

		cookiesData, err := network.GetCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}

		if err = os.WriteFile(path, cookiesData, 0755); err != nil {
			return
		}

		return
	}
}

// saveCode save qr code
func saveCode(url, sel, path string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if url == "" || sel == "" || path == "" {
			return
		}

		chromedp.Navigate(url)
		if err = chromedp.WaitVisible(sel, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		var code []byte
		if err = chromedp.Screenshot(sel, &code, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		if err = os.WriteFile(path, code, 0755); err != nil {
			return
		}

		log.Printf("save qr code at:%s\n", path)
		return
	}
}

// TODO Ref:https://github.com/chromedp/chromedp/issues/484
func getCode1(selector string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 1. 用于存储图片的字节切片
		var code []byte

		template := `
			var img = document.querySelector('%s');
			var c = document.createElement('canvas');
			
			c.height = img.naturalHeight;
			c.width = img.naturalWidth;
			var ctx = c.getContext('2d');
			ctx.drawImage(img, 0, 0, c.width, c.height);
			c.toDataURL('image/png');`

		chromedp.Evaluate(fmt.Sprintf(template, selector), &code)

		// 3. 保存文件
		if err = os.WriteFile("code.png", code, 0755); err != nil {
			log.Printf("%s", err)
			return
		}

		return
	}
}
