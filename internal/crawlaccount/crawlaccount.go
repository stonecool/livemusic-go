package crawlaccount

import (
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
)

type CrawlAccount struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	AccountName  string `json:"account_name"`
	lastURL      string
	cookies      []byte
	instanceAddr string
	state        internal.AccountState
	mu           sync.RWMutex
	msgChan      chan *internal.AsyncMessage // 改用 AsyncMessage
	done         chan struct{}
}

func NewCrawlAccount(m *model.CrawlAccount) *CrawlAccount {
	ca := &CrawlAccount{
		ID:          m.ID,
		Category:    m.Category,
		AccountName: m.AccountName,
		cookies:     m.Cookies,
		lastURL:     m.LastURL,
		msgChan:     make(chan *internal.AsyncMessage),
		done:        make(chan struct{}),
	}
	go ca.processTask()
	return ca
}

func (ca *CrawlAccount) processTask() {
	for {
		select {
		case msg := <-ca.msgChan:
			switch msg.Cmd {
			case internal.CrawlCmd_Login:
				result := ca.handleLogin()
				msg.Data = []byte(fmt.Sprintf("%v", result))
			case internal.CrawlCmd_Crawl:
				result := ca.handleCrawl(msg.Data)
				msg.Data = []byte(fmt.Sprintf("%v", result))
			}
		case <-ca.done:
			return
		}
	}
}

func (ca *CrawlAccount) Close() {
	close(ca.done)
}

func (ca *CrawlAccount) handleLogin() interface{} {
	// 处理登录任务的具体逻辑
	return nil
}

func (ca *CrawlAccount) handleCrawl(payload interface{}) interface{} {
	// 处理爬取任务的具体逻辑
	return nil
}

func (ca *CrawlAccount) Add() error {
	_, ok := config.AccountMap[ca.Category]
	if !ok {
		return fmt.Errorf("account_type:%s not exists", ca.Category)
	}

	data := map[string]interface{}{
		"account_type": ca.Category,
	}

	if _, err := model.AddCrawlAccount(data); err != nil {
		return err
	} else {
		return nil
	}
}

func (ca *CrawlAccount) Get() error {
	//if account, err := model.GetCrawlAccount(ca.ID); err != nil {
	//	return err
	//} else {
	return nil
	//}
}

// FIXME
// func (ca *CrawlAccount) GetAll() ([]*CrawlAccount, error) {
// 	if accounts, err := model.GetCrawlAccountAll(); err != nil {
// 		return nil, err
// 	} else {
// 		var s []*CrawlAccount

// 		for _, account := range accounts {
// 			s = append(s, newCrawlAccount(account))
// 		}

// 		return s, nil
// 	}
// }

func (ca *CrawlAccount) Edit() error {
	if ca.ID == 0 {
		return fmt.Errorf("invalid account id")
	}

	data := map[string]interface{}{
		"account_name":   ca.AccountName,
		"last_login_url": ca.lastURL,
		"cookies":        ca.cookies,
	}

	return model.EditCrawlAccount(ca.ID, data)
}

func (ca *CrawlAccount) Delete() error {
	account, err := model.GetCrawlAccount(ca.ID)
	if err != nil {
		return err
	}

	return model.DeleteCrawlAccount(account)
}

func (ca *CrawlAccount) GetId() int {
	return ca.ID
}

func (ca *CrawlAccount) GetName() string {
	return ca.AccountName
}

func (ca *CrawlAccount) GetCategory() string {
	return ca.Category
}

func (ca *CrawlAccount) GetState() internal.AccountState {
	return ca.state
}

func (ca *CrawlAccount) SetState(state internal.AccountState) {
	ca.state = state
}

func (ca *CrawlAccount) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (ca *CrawlAccount) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (ca *CrawlAccount) GetLoginURL() string {
	return ""
}

func (ca *CrawlAccount) Login() error {
	return nil
}

func (ca *CrawlAccount) GetQRCode([]byte) {
}

func (ca *CrawlAccount) GetQRCodeSelector() string {
	return ""
}

func (ca *CrawlAccount) SaveCookies([]byte) error {
	return nil
}

func (ca *CrawlAccount) GetCookies() []byte {
	return nil
}

func (ca *CrawlAccount) GetLastURL() string {
	return ""
}

func (ca *CrawlAccount) SetLastURL(url string) {
}

func (ca *CrawlAccount) IsAvailable() bool {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	return ca.state == internal.AS_RUNNING
}

func (ca *CrawlAccount) GetMsgChan() chan *internal.AsyncMessage {
	return ca.msgChan
}
