package account

import (
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/client"

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
	msgChan     chan *client.AsyncMessage
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
			case client.CrawlCmd_Login:
				result := acc.handleLogin()
				msg.Data = []byte(fmt.Sprintf("%v", result))
			case client.CrawlCmd_Crawl:
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

//func (ca *Account) Add() error {
//	_, ok := config.caountMap[ca.Category]
//	if !ok {
//		return fmt.Errorf("caount_type:%s not exists", ca.Category)
//	}
//
//	data := map[string]interface{}{
//		"caount_type": ca.Category,
//	}
//
//	if _, err := database.Addcaount(data); err != nil {
//		return err
//	} else {
//		return nil
//	}
//}

func (acc *Account) Get() error {
	//if crawlcaount, err := database.Getcaount(acc.ID); err != nil {
	//	return err
	//} else {
	return nil
	//}
}

// FIXME
// func (ca *Account) GetAll() ([]*Account, error) {
// 	if caounts, err := database.GetcaountAll(); err != nil {
// 		return nil, err
// 	} else {
// 		var s []*Account

// 		for _, crawlcaount := range caounts {
// 			s = append(s, newcaount(crawlcaount))
// 		}

// 		return s, nil
// 	}
// }

//func (ca *Account) Edit() error {
//	if ca.ID == 0 {
//		return fmt.Errorf("invalid crawlcaount id")
//	}
//
//	data := map[string]interface{}{
//		"caount_name":   ca.caountName,
//		"last_login_url": ca.lastURL,
//		"cookies":        ca.cookies,
//	}
//
//	return database.Editcaount(ca.ID, data)
//}
//
//func (ca *Account) Delete() error {
//	crawlcaount, err := database.Getcaount(ca.ID)
//	if err != nil {
//		return err
//	}
//
//	return database.Deletecaount(crawlcaount)
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

func (acc *Account) GetMsgChan() chan *client.AsyncMessage {
	return acc.msgChan
}

func (acc *Account) GetID() int {
	return acc.ID
}
