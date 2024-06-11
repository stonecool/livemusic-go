package crawl

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/model"
	"log"
	"reflect"
)

type Crawl struct {
	ID          int    `json:"id"`
	AccountType string `json:"account_type"`
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	cookies     map[string]interface{}
	State       uint8 `json:"state"`
	config      *internal.AccountConfig
	account     *Crawl
	ch          chan []byte
}

func AddCrawl(crawlType string) (*Crawl, error) {
	data := map[string]interface{}{
		"crawl_type": crawlType,
		"state":      internal.CsNotLoggedIn,
	}

	if m, err := model.AddCrawl(data); err != nil {
		return nil, err
	} else {
		crawl := Crawl{
			ID:          m.ID,
			AccountType: m.CrawlType,
			State:       m.State,
		}

		return &crawl, nil
	}
}

func GetCrawlByID(id int) (*Crawl, error) {
	m, err := model.GetCrawl(id)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	if reflect.ValueOf(*m).IsZero() {
		return &Crawl{}, nil
	}

	crawl := Crawl{
		ID: m.ID,
	}

	return &crawl, nil
}

func (c *Crawl) GetCookies() []byte {
	return nil
}

func (c *Crawl) GetChan() chan []byte {
	return nil
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

// websocket
// 先启动，如果有cookie，尝试自动登录
// 如果自动登录失败，返回“未登录”
// 扫码登录
// getCrawl
func getCrawl(a *Crawl) ICrawl {
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
		crawl = &WxCrawl{
			Crawl: Crawl{
				config:  &cfg,
				account: a,
				ch:      make(chan []byte),
			},
		}
	}

	return crawl
}
