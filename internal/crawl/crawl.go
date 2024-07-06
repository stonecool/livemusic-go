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
	State       uint8  `json:"state"`
	cookies     []byte
	config      *internal.CrawlAccount
	ch          chan *internal.Message
}

var crawlInstances *cache.Memo

func init() {
	crawlInstances = cache.New(getCrawlByID)
}

func AddCrawl(crawlType string) (*Crawl, error) {
	_, ok := internal.CrawlAccountMap[crawlType]
	if !ok {
		// FIXME
		return nil, error(nil)
	}

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
		}

		return &crawl, nil
	}
}

func (c *Crawl) GetId() string {
	return c.AccountId
}

func (c *Crawl) GetName() string {
	return c.AccountName
}

func (c *Crawl) GetState() internal.CrawlState {
	return internal.CrawlState_Ready
}

func (c *Crawl) SetState(state internal.CrawlState) {
	//c.State = state.
}

func (c *Crawl) GetCookies() []byte {
	return nil
}

func (c *Crawl) GetChan() chan *internal.Message {
	return nil
}

func (c *Crawl) SetName(name string) {
	c.AccountName = name
}

func (c *Crawl) CheckLogin() chromedp.ActionFunc {
	return nil
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
				ch:     make(chan *internal.Message),
			},
		}
	}

	go crawl.Start()
	return crawl, nil
}

// GetCrawlByID
func GetCrawlByID(id int) (*Crawl, error) {
	crawl, err := crawlInstances.Get(id)
	if err != nil {
		return nil, err
	} else {
		return crawl.(*Crawl), nil
	}
}

func (c *Crawl) Start() {
	log.Printf("Start crawl:%d\n", c.GetId())

	for {
		select {
		case msg := <-c.GetChan():
			curState := c.GetState()

			switch msg.Cmd {
			case internal.CrawlCmd_Initial:
				if curState != internal.CrawlState_Uninitialized {
					continue
				}

				ret, err := c.Login()
				if err != nil {
					log.Printf("error:%s", err)
					continue
				}

				if ret {
					c.SetState(internal.CrawlState_NotLogged)
				}

			case internal.CrawlCmd_Login:
				if curState != internal.CrawlState_NotLogged {
					log.Printf("state not ready")
					continue
				}

				c.SetState(internal.CrawlState_Ready)

			case internal.CrawlCmd_Crawl:

			default:
				log.Printf("cmd:%v not supportted", msg.Cmd)
			}
		}
	}
}
