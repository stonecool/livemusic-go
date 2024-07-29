package internal

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/config"
)

type Crawl struct {
	Account *CrawlAccount

	state  CrawlState
	config *config.Account
	ch     chan *ClientMessage
}

func (c *Crawl) GetId() int {
	return c.Account.ID
}

func (c *Crawl) GetName() string {
	return c.Account.AccountName
}

func (c *Crawl) GetState() CrawlState {
	return c.state
}

func (c *Crawl) SetState(state CrawlState) {
	c.state = state
}

func (c *Crawl) Login() (bool, error) {
	return false, nil
}

func (c *Crawl) CheckLogin() chromedp.ActionFunc {
	return nil
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

func (c *Crawl) GetCookies() []byte {
	return nil
}

func (c *Crawl) GetChan() chan *ClientMessage {
	return c.ch
}

func (c *Crawl) GetLastLoginURL() string {
	if len(c.Account.lastLoginURL) != 0 {
		return c.Account.lastLoginURL
	}

	return c.config.LastLoginURL
}

func (c *Crawl) SetLastLoginURl(url string) {
	c.Account.lastLoginURL = url
}
