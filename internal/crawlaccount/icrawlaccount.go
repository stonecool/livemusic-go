package crawlaccount

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
)

type ICrawlAccount interface {
	GetId() int

	GetName() string

	GetCategory() string

	GetState() internal.AccountState

	SetState(state internal.AccountState)

	CheckLogin() chromedp.ActionFunc

	WaitLogin() chromedp.ActionFunc

	GetLoginURL() string

	Login() error

	GetQRCode([]byte)

	GetQRCodeSelector() string

	SaveCookies([]byte) error

	GetCookies() []byte

	GetLastURL() string

	SetLastURL(url string)

	GetMsgChan() chan *internal.AsyncMessage
}
