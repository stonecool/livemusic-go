package account

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/message"
)

type IAccount interface {
	Init()
	GetID() int
	GetName() string
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
	GetMsgChan() chan *message.AsyncMessage
}