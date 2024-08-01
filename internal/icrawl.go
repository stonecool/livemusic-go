package internal

import (
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

	//crawl(instance *Instance) error
	//
	//GetLoginRequestData() []byte
	//
	//LoginRequestCallback(request *colly.Request) error
	//
	//LoginResponseCallback(response *colly.Response) error
}
