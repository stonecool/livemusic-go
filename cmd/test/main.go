package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

func createContextWithTimeout() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	allocatorCtx, cancel := chromedp.NewRemoteAllocator(ctx, "")
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	return browserCtx
}

func main() {

}
