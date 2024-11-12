package crawlaccount

import "github.com/chromedp/chromedp"

type ICrawlAccount interface {
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

	GetLastLoginURL() string

	SetLastLoginURL(url string)

	//GoCrawl(Callback) chromedp.ActionFunc

	Login() error
}
