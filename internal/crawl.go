package internal

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/cache"
	crawl2 "github.com/stonecool/livemusic-go/internal/crawl"
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
	config      *CrawlAccount
	ch          chan *Message
}

var crawlInstances *cache.Memo

func init() {
	crawlInstances = cache.New(getCrawlByID)
}

func (crawl *Crawl) init(m *model.Crawl) {
	if m == nil || reflect.ValueOf(m).IsZero() {
		return
	}

	crawl.ID = m.ID
	crawl.CrawlType = m.CrawlType
	crawl.AccountId = m.AccountId
	crawl.AccountName = m.AccountName
	crawl.cookies = m.Cookies
}

// AddCrawl
func AddCrawl(crawlType string) (*Crawl, error) {
	_, ok := CrawlAccountMap[crawlType]
	if !ok {
		return nil, error(nil)
	}

	data := map[string]interface{}{
		"crawl_type": crawlType,
		"state":      uint8(0),
	}

	if m, err := model.AddCrawl(data); err != nil {
		return nil, err
	} else {
		crawl := Crawl{}
		crawl.init(m)
		return &crawl, nil
	}
}

func (crawl *Crawl) GetId() string {
	return crawl.AccountId
}

func (crawl *Crawl) GetName() string {
	return crawl.AccountName
}

func (crawl *Crawl) GetState() CrawlState {
	return CrawlState_Ready
}

func (crawl *Crawl) SetState(state CrawlState) {
	//crawl.State = state.
}

func (crawl *Crawl) GetCookies() []byte {
	return nil
}

func (crawl *Crawl) GetChan() chan *Message {
	return nil
}

func (crawl *Crawl) SetName(name string) {
	crawl.AccountName = name
}

func (crawl *Crawl) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (crawl *Crawl) Login() (bool, error) {
	return false, nil
}

func (crawl *Crawl) GetLoginURL() string {
	return crawl.config.LoginURL
}

func (crawl *Crawl) GetQRCode(data []byte) {

}

func (crawl *Crawl) GetQRCodeSelector() string {
	return ""
}

func (crawl *Crawl) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (crawl *Crawl) SaveCookies([]byte) {
}

func (crawl *Crawl) GetLoginSelector() string {
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

	cfg, ok := CrawlAccountMap[m.CrawlType]
	if !ok {
		return nil, nil
	}

	var crawl crawl2.ICrawl
	switch m.CrawlType {
	case "wx":
		crawl = &crawl2.WxCrawl{
			Crawl: Crawl{
				ID:     m.ID,
				config: &cfg,
				ch:     make(chan *Message),
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

func (crawl *Crawl) Start() {
	log.Printf("Start crawl:%d\n", crawl.GetId())

	for {
		select {
		case msg := <-crawl.GetChan():
			curState := crawl.GetState()

			switch msg.Cmd {
			case CrawlCmd_Initial:
				if curState != CrawlState_Uninitialized {
					continue
				}

				ret, err := crawl.Login()
				if err != nil {
					log.Printf("error:%s", err)
					continue
				}

				if ret {
					crawl.SetState(CrawlState_NotLogged)
				}

			case CrawlCmd_Login:
				if curState != CrawlState_NotLogged {
					log.Printf("state not ready")
					continue
				}

				crawl.SetState(CrawlState_Ready)

			case CrawlCmd_Crawl:

			default:
				log.Printf("cmd:%v not supportted", msg.Cmd)
			}
		}
	}
}
