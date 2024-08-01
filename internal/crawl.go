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

func (c *Crawl) SaveCookies(cookies []byte) error {
	c.Account.cookies = cookies
	return c.Account.Edit()
}

func (c *Crawl) GetCookies() []byte {
	return c.Account.cookies
}

func (c *Crawl) GetChan() chan *ClientMessage {
	return c.ch
}

func (c *Crawl) GetLastLoginURL() string {
	return c.config.LastLoginURL
}

func (c *Crawl) SetLastLoginURL(string) {
}
