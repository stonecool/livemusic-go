package crawl

import (
	"context"
	"github.com/chromedp/chromedp"
)

type ICrawl interface {
	GetId() int

	GetName() string

	GetState() CrawlState

	SetState(state CrawlState)

	CheckLogin() chromedp.ActionFunc

	GetLoginURL() string

	GetQRCode([]byte)

	GetQRCodeSelector() string

	WaitLogin() chromedp.ActionFunc

	SaveCookies([]byte) error

	GetCookies() []byte

	GetChan() chan *ClientMessage

	GetLastLoginURL() string

	SetLastLoginURL(url string)

	GoCrawl(Callback) chromedp.ActionFunc

	SetContext(ctx context.Context)

	GetContext() context.Context

	Login() error

	callback(map[string]interface{})
}
