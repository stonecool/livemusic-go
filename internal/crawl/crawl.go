package crawl

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"reflect"
)

type Crawl struct {
	config  *internal.AccountConfig
	account *Account
	ch      chan []byte
}

func (c *Crawl) GetId() string {
	return c.account.AccountId
}

func (c *Crawl) setId(id string) {
	c.account.AccountId = id
}

func (c *Crawl) GetName() string {
	return c.account.AccountName
}

func (c *Crawl) SetName(name string) {
	c.account.AccountName = name
}

func (c *Crawl) CheckLogin() (bool, error) {
	return false, nil
}

func (c *Crawl) Login() (bool, error) {
	return false, nil
}

func (c *Crawl) GetLoginURL() string {
	return c.config.LoginURL
}

func (c *Crawl) GetQRCode(data []byte) {

}

func (c *Crawl) GetQRCodeSelector() string {
	return ""
}

func (c *Crawl) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (c *Crawl) SaveCookies([]byte) {
}

func (c *Crawl) GetLoginSelector() string {
	return ""
}

func (c *Crawl) GetCookies() []byte {
	return c.account.GetCookies()
}

func (c *Crawl) GetChan() chan []byte {
	return nil
}

// websocket
// 先启动，如果有cookie，尝试自动登录
// 如果自动登录失败，返回“未登录”
// 扫码登录
// GetCrawl
func GetCrawl(a *Account) ICrawl {
	// FIXME
	if a == nil || reflect.ValueOf(a).IsZero() {
		return nil
	}

	cfg, ok := internal.AccountConfigMap[a.AccountType]
	if !ok {
		return nil
	}

	var crawl ICrawl
	switch a.AccountType {
	case "WxPublicAccount":
		crawl = &WxPublicAccountCrawl{
			Crawl: Crawl{
				config:  &cfg,
				account: a,
				ch:      make(chan []byte),
			},
		}
	}

	return crawl
}
