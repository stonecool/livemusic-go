package crawlaccount

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"sync"

	"github.com/stonecool/livemusic-go/internal"
)

type Account struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	AccountName string `json:"account_name"`
	lastURL     string
	cookies     []byte
	InstanceID  int
	State       internal.AccountState
	mu          sync.RWMutex
	msgChan     chan *internal.AsyncMessage
	done        chan struct{}
}

func (acc *Account) Init() {
	go acc.processTask()
}

func (acc *Account) processTask() {
	for {
		select {
		case msg := <-acc.msgChan:
			switch msg.Cmd {
			case internal.CrawlCmd_Login:
				result := acc.handleLogin()
				msg.Data = []byte(fmt.Sprintf("%v", result))
			case internal.CrawlCmd_Crawl:
				result := acc.handleCrawl(msg.Data)
				msg.Data = []byte(fmt.Sprintf("%v", result))
			}
		case <-acc.done:
			return
		}
	}
}
func (acc *Account) Close() {
	close(acc.done)
}

func (acc *Account) handleLogin() interface{} {
	// 处理登录任务的具体逻辑
	return nil
}

func (acc *Account) handleCrawl(payload interface{}) interface{} {
	// 处理爬取任务的具体逻辑
	return nil
}

//func (acc *Account) Add() error {
//	_, ok := config.AccountMap[ca.Category]
//	if !ok {
//		return fmt.Errorf("account_type:%s not exists", ca.Category)
//	}
//
//	data := map[string]interface{}{
//		"account_type": ca.Category,
//	}
//
//	if _, err := model.AddAccount(data); err != nil {
//		return err
//	} else {
//		return nil
//	}
//}

func (acc *Account) Get() error {
	//if crawlaccount, err := model.GetAccount(ca.ID); err != nil {
	//	return err
	//} else {
	return nil
	//}
}

// FIXME
// func (acc *Account) GetAll() ([]*Account, error) {
// 	if accounts, err := model.GetAccountAll(); err != nil {
// 		return nil, err
// 	} else {
// 		var s []*Account

// 		for _, crawlaccount := range accounts {
// 			s = append(s, newAccount(crawlaccount))
// 		}

// 		return s, nil
// 	}
// }

//func (acc *Account) Edit() error {
//	if ca.ID == 0 {
//		return fmt.Errorf("invalid crawlaccount id")
//	}
//
//	data := map[string]interface{}{
//		"account_name":   ca.AccountName,
//		"last_login_url": ca.lastURL,
//		"cookies":        ca.cookies,
//	}
//
//	return model.EditAccount(ca.ID, data)
//}
//
//func (acc *Account) Delete() error {
//	crawlaccount, err := model.GetAccount(ca.ID)
//	if err != nil {
//		return err
//	}
//
//	return model.DeleteAccount(crawlaccount)
//}

func (acc *Account) GetId() int {
	return acc.ID
}

func (acc *Account) GetName() string {
	return acc.AccountName
}

func (acc *Account) GetCategory() string {
	return acc.Category
}

func (acc *Account) GetState() internal.AccountState {
	acc.mu.RLock()
	defer acc.mu.RUnlock()

	return acc.State
}

func (acc *Account) SetState(state internal.AccountState) {
	acc.mu.Lock()
	defer acc.mu.Lock()

	acc.State = state
}

func (acc *Account) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (acc *Account) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (acc *Account) GetLoginURL() string {
	return ""
}

func (acc *Account) Login() error {
	return nil
}

func (acc *Account) GetQRCode([]byte) {
}

func (acc *Account) GetQRCodeSelector() string {
	return ""
}

func (acc *Account) SaveCookies([]byte) error {
	return nil
}

func (acc *Account) GetCookies() []byte {
	return nil
}

func (acc *Account) GetLastURL() string {
	return ""
}

func (acc *Account) SetLastURL(url string) {
}

func (acc *Account) IsAvailable() bool {
	acc.mu.Lock()
	defer acc.mu.Unlock()

	return acc.State == internal.AS_RUNNING
}

func (acc *Account) GetMsgChan() chan *internal.AsyncMessage {
	return acc.msgChan
}

func (acc *Account) GetID() int {
	return acc.ID
}