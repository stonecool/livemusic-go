package crawl

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/model"
	"log"
	"reflect"
)

type Crawl struct {
	ID          int    `json:"id"`
	CrawlType   string `json:"crawl_type"`
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	cookies     []byte
	State       uint8 `json:"state"`
	config      *internal.CrawlAccount
	ch          chan []byte
}

var crawlInstances *cache.Memo

func init() {
	crawlInstances = cache.New(getCrawlByID)
}

func AddCrawl(crawlType string) (*Crawl, error) {
	data := map[string]interface{}{
		"crawl_type": crawlType,
		"state":      uint8(0),
	}

	if m, err := model.AddCrawl(data); err != nil {
		return nil, err
	} else {
		crawl := Crawl{
			ID:        m.ID,
			CrawlType: m.CrawlType,
			State:     m.State,
		}

		return &crawl, nil
	}
}

func GetCrawlByID(id int) (*Crawl, error) {
	crawl, err := crawlInstances.Get(id)
	if err != nil {
		return nil, err
	} else {
		return crawl.(*Crawl), nil
	}
}

func (c *Crawl) GetCookies() []byte {
	return nil
}

func (c *Crawl) GetChan() chan []byte {
	return nil
}

func (c *Crawl) GetId() string {
	return c.AccountId
}

func (c *Crawl) GetName() string {
	return c.AccountName
}

func (c *Crawl) SetName(name string) {
	c.AccountName = name
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

func (c *Crawl) Start() {
}

// websocket
// 先启动，如果有cookie，尝试自动登录
// 如果自动登录失败，返回“未登录”
// 扫码登录
// getCrawlByID
func getCrawlByID(id int) (interface{}, error) {
	m, err := model.GetCrawlByID(id)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	if reflect.ValueOf(*m).IsZero() {
		return &Crawl{}, nil
	}

	cfg, ok := internal.CrawlAccountMap[m.CrawlType]
	if !ok {
		return nil, nil
	}

	var crawl ICrawl
	switch m.CrawlType {
	case "wx":
		crawl = &WxCrawl{
			Crawl: Crawl{
				ID:     m.ID,
				config: &cfg,
				ch:     make(chan []byte),
			},
		}
	}

	go crawl.Start()
	return crawl, nil
}
