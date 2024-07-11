package internal

import (
	"github.com/stonecool/livemusic-go/internal/model"
	"reflect"
)

type CrawlAccount struct {
	ID          int    `json:"id"`
	CrawlType   string `json:"crawl_type"`
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	State       uint8  `json:"state"`
	cookies     []byte
	config      *Account
	ch          chan *Message
}

func (crawl *CrawlAccount) init(m *model.CrawlAccount) {
	if m == nil || reflect.ValueOf(m).IsZero() {
		return
	}

	crawl.ID = m.ID
	crawl.CrawlType = m.CrawlType
	crawl.AccountId = m.AccountId
	crawl.AccountName = m.AccountName
	crawl.cookies = m.Cookies
}

// AddCrawlAccount
func AddCrawlAccount(crawlType string) (*CrawlAccount, error) {
	_, ok := CrawlAccountMap[crawlType]
	if !ok {
		return nil, error(nil)
	}

	data := map[string]interface{}{
		"crawl_type": crawlType,
		"state":      uint8(0),
	}

	if m, err := model.AddCrawlAccount(data); err != nil {
		return nil, err
	} else {
		crawl := CrawlAccount{}
		crawl.init(m)
		return &crawl, nil
	}
}

// getCrawlByID
func getCrawlByID(id int) (*CrawlAccount, error) {
	//m, err := model.GetCrawlAccount(id)
	//if err != nil {
	//	log.Printf("error: %s", err)
	//	return nil, err
	//}
	//
	//if reflect.ValueOf(*m).IsZero() {
	//	return &Account{}, nil
	//}
	//
	//cfg, ok := CrawlAccountMap[m.CrawlType]
	//if !ok {
	//	return nil, nil
	//}

	//var crawl crawl2.ICrawl
	//switch m.CrawlType {
	//case "wx":
	//	crawl = &crawl2.WxCrawl{
	//		Account: Account{
	//			ID:     m.ID,
	//			config: &cfg,
	//			ch:     make(chan *Message),
	//		},
	//	}
	//}

	//go crawl.Start()
	return nil, nil
}

// GetCrawlByID
func GetCrawlByID(id int) (*CrawlAccount, error) {
	crawl, ok := crawlInGoroutine[id]
	if ok {
		return crawl, nil
	}

	crawl, err := getCrawlByID(id)
	if err != nil {
		return nil, err
	} else {
		return crawl, nil
	}
}
