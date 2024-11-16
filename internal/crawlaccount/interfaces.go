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

type IRepository interface {
	Create(account *Account) error
	Get(id int) (*Account, error)
	Update(account *Account) error
	Delete(id int) error
	GetAll() ([]*Account, error)
	FindByCategory(category string) ([]*Account, error)
	FindByInstance(instanceID int) ([]*Account, error)
}

type IFactory interface {
	CreateAccount(category string) (*Account, error)
}
