package crawlaccount

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"sync"

	"github.com/stonecool/livemusic-go/internal"
)

type CrawlAccount struct {
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

func (ca *CrawlAccount) Init() {
	go ca.processTask()
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

//func (ca *CrawlAccount) Add() error {
//	_, ok := config.caountMap[ca.Category]
//	if !ok {
//		return fmt.Errorf("caount_type:%s not exists", ca.Category)
//	}
//
//	data := map[string]interface{}{
//		"caount_type": ca.Category,
//	}
//
//	if _, err := model.Addcaount(data); err != nil {
//		return err
//	} else {
//		return nil
//	}
//}

func (ca *CrawlAccount) Get() error {
	//if crawlcaount, err := model.Getcaount(ca.ID); err != nil {
	//	return err
	//} else {
	return nil
	//}
}

// FIXME
// func (ca *CrawlAccount) GetAll() ([]*CrawlAccount, error) {
// 	if caounts, err := model.GetcaountAll(); err != nil {
// 		return nil, err
// 	} else {
// 		var s []*CrawlAccount

// 		for _, crawlcaount := range caounts {
// 			s = append(s, newcaount(crawlcaount))
// 		}

// 		return s, nil
// 	}
// }

//func (ca *CrawlAccount) Edit() error {
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
//	return model.Editcaount(ca.ID, data)
//}
//
//func (ca *CrawlAccount) Delete() error {
//	crawlcaount, err := model.Getcaount(ca.ID)
//	if err != nil {
//		return err
//	}
//
//	return model.Deletecaount(crawlcaount)
//}

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
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	return ca.State
}

func (ca *CrawlAccount) SetState(state internal.AccountState) {
	ca.mu.Lock()
	defer ca.mu.Lock()

	ca.State = state
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

	return ca.State == internal.AS_RUNNING
}

func (ca *CrawlAccount) GetMsgChan() chan *internal.AsyncMessage {
	return ca.msgChan
}

func (ca *CrawlAccount) GetID() int {
	return ca.ID
}
