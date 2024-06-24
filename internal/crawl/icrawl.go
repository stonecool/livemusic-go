package crawl

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
)

type ICrawl interface {
	GetId() string

	GetName() string

	GetState() internal.CrawlState

	SetState(state internal.CrawlState)

	Login() (bool, error)

	CheckLogin() chromedp.ActionFunc

	GetLoginURL() string

	GetQRCode([]byte)

	GetQRCodeSelector() string

	WaitLogin() chromedp.ActionFunc

	SaveCookies([]byte)

	GetCookies() []byte

	GetChan() chan *internal.Message

	Start()

	//crawl(instance *Instance) error
	//
	//GetLoginRequestData() []byte
	//
	//LoginRequestCallback(request *colly.Request) error
	//
	//LoginResponseCallback(response *colly.Response) error
}
