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
		}

		return &crawl, nil
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

func (c *Crawl) Start() {
	log.Printf("Start crawl:%d\n", c.GetId())

	//for {
	//	select {
	//	case cmd := <-c.GetChan():
	//		switch cmd.cmd {
	//		case CmdReady:
	//			if crawl.GetState() != StateInitial {
	//				log.Printf("state not initial")
	//				continue
	//			}
	//
	//			ret, err := crawl.Login()
	//			if err != nil {
	//				log.Printf("error:%s", err)
	//				continue
	//			}
	//
	//			if ret {
	//				crawl.SetState(StateReady)
	//			}
	//		case CmdRun:
	//			if crawl.GetState() != StateReady {
	//				log.Printf("state not ready")
	//				continue
	//			}
	//
	//			crawl.SetState(StateRunning)
	//		case CmdSuspend:
	//			if crawl.GetState() != StateRunning {
	//				log.Printf("state not running")
	//				continue
	//			}
	//
	//			crawl.SetState(StateReady)
	//		case CmdCrawl:
	//			if crawl.GetState() != StateRunning {
	//				log.Printf("state not running")
	//				continue
	//			}
	//
	//			err := crawl.crawl(cmd.instance)
	//			if err != nil {
	//				log.Printf("error:%s", err)
	//				// TODO
	//			}
	//		}
	//	}
	//}
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
				ch:     make(chan []byte),
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
