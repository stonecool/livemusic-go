package crawl

import "github.com/chromedp/chromedp"

type ICrawl interface {
	GetId() string

	GetName() string

	//GetChan() chan CmdRequest
	//
	//GetState() State
	//
	//SetState(state State)

	Login() (bool, error)

	CheckLogin() chromedp.ActionFunc

	GetLoginURL() string

	GetQRCode([]byte)

	GetQRCodeSelector() string

	WaitLogin() chromedp.ActionFunc

	SaveCookies([]byte)

	GetCookies() []byte

	GetChan() chan []byte

	Start()

	//crawl(instance *Instance) error
	//
	//GetLoginRequestData() []byte
	//
	//LoginRequestCallback(request *colly.Request) error
	//
	//LoginResponseCallback(response *colly.Response) error
}
