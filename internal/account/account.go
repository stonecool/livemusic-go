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
	State        state
	mu           sync.RWMutex
	msgChan      chan *message.AsyncMessage
	done         chan struct{}
	stateManager stateManager
}

func (acc *account) processTask() {
	for {
		select {
		case msg := <-acc.msgChan:
			state := acc.getState()
			var err error

			ret := acc.handleCommand(state, msg)
			err = ret.(error)

			if err != nil {
				newState := acc.stateManager.getErrorState(state)
				if acc.stateManager.isValidTransition(state, newState) {
					acc.SetState(newState)
				}
			} else {
				newState := acc.stateManager.getNextState(state, msg.Cmd)
				if acc.stateManager.isValidTransition(state, newState) {
					acc.SetState(newState)
				}
			}

			if msg.Result != nil {
				msg.Result <- err
				close(msg.Result)
			}

		case <-acc.done:
			return
		}
	}
}

func (acc *account) handleCommand(currentState state, msg *message.AsyncMessage) interface{} {
	switch currentState {
	case stateNew:
		if msg.Cmd != message.CrawlCmd_Initial {
			return fmt.Errorf("invalid command:%v for new state", msg.Cmd)
		}
		go acc.processTask()
		return nil

	case stateInitialized:
		if msg.Cmd != message.CrawlCmd_Login {
			return fmt.Errorf("invalid command:%v for initialized state", msg.Cmd)
		}
		return acc.handleLogin()

	case stateNotLoggedIn:
		if msg.Cmd != message.CrawlCmd_Login {
			return fmt.Errorf("invalid command:%v for not logged in state", msg.Cmd)
		}
		return acc.handleLogin()

	case stateReady:
		switch msg.Cmd {
		case message.CrawlCmd_Crawl:
			return acc.handleCrawl(msg.Data)
		case message.CrawlCmd_Login:
			return acc.handleLogin()
		default:
			return fmt.Errorf("invalid command:%v for ready state", msg.Cmd)
		}

	case stateRunning:
		return fmt.Errorf("account is busy")

	case stateTerminated:
		return fmt.Errorf("account is terminated")

	default:
		return nil
	}
}

func (acc *account) Close() {
	close(acc.done)
}

func (acc *account) handleLogin() interface{} {
	// 处理登录任务的具体逻辑
	return nil
}

func (acc *account) handleCrawl(payload interface{}) interface{} {
	// 处理爬取任务的具体逻辑
	return payload
}

func (acc *account) Get() error {
	//if crawlcaount, err := database.Getcaount(acc.ID); err != nil {
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

func (acc *account) GetName() string {
	return acc.Name
}

func (acc *account) getState() state {
	acc.mu.RLock()
	defer acc.mu.RUnlock()

	return acc.State
}

func (acc *account) SetState(s state) {
	acc.mu.Lock()
	defer acc.mu.Unlock()

	acc.State = s
}

func (acc *account) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (acc *account) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (acc *account) GetLoginURL() string {
	return ""
}

func (acc *account) Login() error {
	return nil
}

func (acc *account) GetQRCode([]byte) {
}

func (acc *account) GetQRCodeSelector() string {
	return ""
}

func (acc *account) SaveCookies([]byte) error {
	return nil
}

func (acc *account) GetCookies() []byte {
	return nil
}

func (acc *account) GetLastURL() string {
	return ""
}

func (acc *account) SetLastURL(url string) {
}

func (acc *account) IsAvailable() bool {
	acc.mu.Lock()
	defer acc.mu.Unlock()

	return acc.State == stateInitialized
}

func (acc *account) GetMsgChan() chan *message.AsyncMessage {
	return acc.msgChan
}

func (acc *account) GetID() int {
	return acc.ID
}
