package crawlaccount

import (
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
)

type Task struct {
	Type    string
	Payload interface{}
	Result  chan interface{}
}

type CrawlAccount struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	AccountName  string `json:"account_name"`
	lastURL      string
	cookies      []byte
	instanceAddr string
	state        internal.AccountState
	mu           sync.RWMutex
	taskChan     chan *Task
	done         chan struct{}
}

func NewCrawlAccount(m *model.CrawlAccount) *CrawlAccount {
	ca := &CrawlAccount{
		ID:          m.ID,
		Category:    m.Category,
		AccountName: m.AccountName,
		cookies:     m.Cookies,
		lastURL:     m.LastURL,
		taskChan:    make(chan *Task),
		done:        make(chan struct{}),
	}
	go ca.processTask()
	return ca
}

func (ca *CrawlAccount) processTask() {
	for {
		select {
		case task := <-ca.taskChan:
			switch task.Type {
			case "login":
				result := ca.handleLogin()
				task.Result <- result
			case "crawl":
				result := ca.handleCrawl(task.Payload)
				task.Result <- result
				// 可以添加更多任务类型
			}
		case <-ca.done:
			return
		}
	}
}

// SendTask 发送任务到账号处理
func (ca *CrawlAccount) SendTask(taskType string, payload interface{}) (interface{}, error) {
	result := make(chan interface{})
	task := &Task{
		Type:    taskType,
		Payload: payload,
		Result:  result,
	}

	select {
	case ca.taskChan <- task:
		return <-result, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("send task timeout")
	}
}

// Close 关闭账号的任务处理
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

// GetTaskChan 返回任务 channel，用于外部发送任务
func (ca *CrawlAccount) GetTaskChan() chan<- *Task {
	return ca.taskChan
}
