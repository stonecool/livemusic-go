package account

import (
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/message"
)

type account struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	lastURL      string
	cookies      []byte
	instanceID   int
	State        accountState
	mu           sync.RWMutex
	msgChan      chan *message.AsyncMessage
	done         chan struct{}
	stateManager stateManager
}

func (a *account) processTask() {
	for {
		select {
		case msg := <-a.msgChan:
			state := a.getState()
			var err error

			ret := a.handleCommand(state, msg)
			err = ret.(error)

			if err != nil {
				newState := a.stateManager.getErrorState(state)
				if a.stateManager.isValidTransition(state, newState) {
					a.SetState(newState)
				}
			} else {
				newState := a.stateManager.getNextState(state, msg.Cmd)
				if a.stateManager.isValidTransition(state, newState) {
					a.SetState(newState)
				}
			}

			if msg.Result != nil {
				msg.Result <- err
				close(msg.Result)
			}

		case <-a.done:
			return
		}
	}
}

func (a *account) handleCommand(currentState accountState, msg *message.AsyncMessage) interface{} {
	switch currentState {
	case stateInitialized:
		if msg.Cmd != message.CrawlCmd_Login {
			return fmt.Errorf("invalid command:%v for initialized accountState", msg.Cmd)
		}
		return a.handleLogin()

	case stateNotLoggedIn:
		if msg.Cmd != message.CrawlCmd_Login {
			return fmt.Errorf("invalid command:%v for not logged in accountState", msg.Cmd)
		}
		return a.handleLogin()

	case stateReady:
		switch msg.Cmd {
		case message.CrawlCmd_Crawl:
			return a.handleCrawl(msg.Data)
		case message.CrawlCmd_Login:
			return a.handleLogin()
		default:
			return fmt.Errorf("invalid command:%v for ready accountState", msg.Cmd)
		}

	case stateRunning:
		return fmt.Errorf("account is busy")

	case stateTerminated:
		return fmt.Errorf("account is terminated")

	default:
		return nil
	}
}

func (a *account) Close() {
	close(a.done)
}

func (a *account) handleLogin() interface{} {
	// 处理登录任务的具体逻辑
	return nil
}

func (a *account) handleCrawl(payload interface{}) interface{} {
	// 处理爬取任务的具体逻辑
	return payload
}

func (a *account) Get() error {
	//if crawlcaount, err := database.Getcaount(a.ID); err != nil {
	//	return err
	//} else {
	return nil
	//}
}

//func (ca *account) Edit() error {
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
//func (ca *account) Delete() error {
//	crawlcaount, err := database.Getcaount(ca.ID)
//	if err != nil {
//		return err
//	}
//
//	return database.Deletecaount(crawlcaount)
//}

func (a *account) GetName() string {
	return a.Name
}

func (a *account) getState() accountState {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.State
}

func (a *account) SetState(s accountState) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.State = s
}

func (a *account) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (a *account) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (a *account) GetLoginURL() string {
	return ""
}

func (a *account) Login() error {
	return nil
}

func (a *account) GetQRCode([]byte) {
}

func (a *account) GetQRCodeSelector() string {
	return ""
}

func (a *account) SaveCookies([]byte) error {
	return nil
}

func (a *account) GetCookies() []byte {
	return nil
}

func (a *account) GetLastURL() string {
	return ""
}

func (a *account) SetLastURL(url string) {
}

func (a *account) IsAvailable() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.State == stateInitialized
}

func (a *account) GetMsgChan() chan *message.AsyncMessage {
	return a.msgChan
}

func (a *account) GetID() int {
	return a.ID
}

func (a *account) GetCategory() string {
	return a.Category
}
