package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func main() {
	// 创建一个新的 Chrome 实例
	ctx, cancel := chromedp.NewExecAllocator(
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

	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()
	// 设置超时
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 20*time.Second)
	defer cancelTimeout()

	// 执行任务
	var res string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(`https://www.google.com`),
		//chromedp.Text(`body`, &res, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(30 * time.Second)
	log.Println("Body text:", res)
}
