package types

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/message"
)

type Account interface {
	GetID() int
	GetName() string
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
	GetCategory() string
}
