package crawl

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/config"
)

type Callback func(map[string]interface{})

type Crawl struct {
	state   internal.CrawlState
	config  *config.Account
	ch      chan *internal.ClientMessage
	context context.Context
}

func (c *Crawl) GetId() int {
	return 0
}

func (c *Crawl) GetName() string {
	return ""
}

func (c *Crawl) GetState() internal.CrawlState {
	return c.state
}

func (c *Crawl) SetState(state internal.CrawlState) {
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
	return nil
}

func (c *Crawl) GetCookies() []byte {
	return nil
}

func (c *Crawl) GetChan() chan *internal.ClientMessage {
	return c.ch
}

func (c *Crawl) GetLastLoginURL() string {
	return c.config.LastLoginURL
}

func (c *Crawl) SetLastLoginURL(string) {
}

func (c *Crawl) GoCrawl() chromedp.ActionFunc {
	return nil
}

func (c *Crawl) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *Crawl) GetContext() context.Context {
	return c.context
}

func (c *Crawl) Login() error {
	return nil
}
