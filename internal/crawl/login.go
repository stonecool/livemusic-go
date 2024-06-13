package crawl

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"time"
)

// QRCodeLogin
func QRCodeLogin(iCrawl ICrawl) error {
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
		getQRCode(iCrawl),
		iCrawl.WaitLogin(),
		iCrawl.CheckLogin(),
		setCookies(iCrawl),
		// TODO if every err should stop?
		chromedp.Stop(),
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// getQRCode get qr code
func getQRCode(iCrawl ICrawl) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		chromedp.Navigate(iCrawl.GetLoginURL())
		if err = chromedp.WaitVisible(iCrawl.GetQRCodeSelector(), chromedp.ByID).Do(ctx); err != nil {
			return
		}

		var code []byte
		if err = chromedp.Screenshot(iCrawl.GetQRCodeSelector(), &code, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		iCrawl.GetQRCode(code)
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

// checkLogin
func checkLogin(iCrawl ICrawl) error {
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
		getCookies(iCrawl),
		iCrawl.WaitLogin(),
		iCrawl.CheckLogin(),
		// TODO if every err should stop?
		chromedp.Stop(),
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
