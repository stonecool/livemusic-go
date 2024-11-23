package account

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/message"
)

type IAccount interface {
	GetID() int
	GetName() string
	getState() state
	SetState(state)
	CheckLogin() chromedp.ActionFunc
	WaitLogin() chromedp.ActionFunc
	GetLoginURL() string
	Login() error
	GetQRCode([]byte)
	GetQRCodeSelector() string
	SaveCookies([]byte) error
	GetCookies() []byte
	GetLastURL() string
	SetLastURL(string)
	GetMsgChan() chan *message.AsyncMessage
	IsAvailable() bool
}
