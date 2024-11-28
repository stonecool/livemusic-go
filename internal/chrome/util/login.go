package util

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account"
	"go.uber.org/zap"
)

// QRCodeLogin
func QRCodeLogin() error {
	//ctx, _ := chromedp.NewExecAllocator(
	//	context.Background(),
	//
	//	append(
	//		chromedp.DefaultExecAllocatorOptions[:],
	//		//chromedp.NoDefaultBrowserCheck,
	//		chromedp.Flag("headless", false),
	//		//chromedp.Flag("hide-scrollbars", false),
	//		//chromedp.Flag("mute-audio", false),
	//		//chromedp.Flag("ignore-certificate-errors", true),
	//		//chromedp.Flag("disable-web-security", true),
	//		//chromedp.Flag("disable-gpu", false),
	//		//chromedp.NoFirstRun,
	//		//chromedp.Flag("enable-automation", false),
	//		//chromedp.Flag("disable-extensions", false),
	//	)...,
	//)
	return nil
}

// GetQRCode get qr code
func GetQRCode(account account.IAccount) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if err := chromedp.Navigate(account.GetLoginURL()).Do(ctx); err != nil {
			internal.Logger.Error("failed to navigate to login URL",
				zap.Error(err),
				zap.Int("account", account.GetID()))
			return err
		}

		if err = chromedp.WaitVisible(account.GetQRCodeSelector(), chromedp.ByID).Do(ctx); err != nil {
			internal.Logger.Error("failed to wait for QR code",
				zap.Error(err),
				zap.Int("account", account.GetID()))
			return err
		}

		var code []byte
		if err = chromedp.Screenshot(account.GetQRCodeSelector(), &code, chromedp.ByID).Do(ctx); err != nil {
			internal.Logger.Error("failed to take screenshot of QR code",
				zap.Error(err),
				zap.Int("account", account.GetID()))
			return err
		}

		if err = os.WriteFile("code.png", code, 0755); err != nil {
			internal.Logger.Error("failed to write QR code file",
				zap.Error(err),
				zap.Int("account", account.GetID()))
			return err
		}

		account.GetQRCode(code)
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
