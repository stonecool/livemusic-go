package account

import (
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/message"
)

type Account struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	AccountName  string `json:"account_name"`
	lastURL      string
	cookies      []byte
	InstanceID   int
	mu           sync.RWMutex
	msgChan      chan *message.AsyncMessage
	done         chan struct{}
	curState     state
	stateManager stateManager
}

func NewAccount(category string) *Account {
	acc := &Account{
		stateManager: selectStateManager(category),
	}
	return acc
}

func (act *Account) Init() {
	go act.processTask()
}

func (act *Account) processTask() {
	for {
		select {
		case msg := <-act.msgChan:
			currentState := act.getState()
			var err error

			// 处理命令
			err = act.handleCommand(currentState, msg)

			// 状态转换
			if err != nil {
				newState := act.stateManager.getErrorState(currentState)
				if act.stateManager.isValidTransition(currentState, newState) {
					act.SetState(newState)
				}
			} else {
				newState := act.stateManager.getNextState(currentState, msg.Cmd)
				if act.stateManager.isValidTransition(currentState, newState) {
					act.SetState(newState)
				}
			}

			// 通过 Result channel 返回结果
			if msg.Result != nil {
				msg.Result <- err
				close(msg.Result)
			}

		case <-act.done:
			return
		}
	}
}

func (act *Account) handleCommand(currentState state, msg *message.AsyncMessage) error {
	// 将原来 switch 中的命令处理逻辑移到这里
	switch currentState {
	case stateNew:
		return fmt.Errorf("invalid command for new state: %v", msg.Cmd)
		// ... 其他状态的命令处理
	}
	return nil
}

func (act *Account) Close() {
	close(act.done)
}

func (act *Account) handleLogin() interface{} {
	// 处理登录任务的具体逻辑
	return nil
}

func (act *Account) handleCrawl(payload interface{}) interface{} {
	// 处理爬取任务的具体逻辑
	return payload
}

func (act *Account) Get() error {
	//if crawlcaount, err := database.Getcaount(act.ID); err != nil {
	//	return err
	//} else {
	return nil
	//}
}

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

func (act *Account) GetName() string {
	return act.AccountName
}

func (act *Account) getState() state {
	act.mu.RLock()
	defer act.mu.RUnlock()

	return act.curState
}

func (act *Account) SetState(s state) {
	act.mu.Lock()
	defer act.mu.Unlock()

	act.curState = s
}

func (act *Account) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (act *Account) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (act *Account) GetLoginURL() string {
	return ""
}

func (act *Account) Login() error {
	return nil
}

func (act *Account) GetQRCode([]byte) {
}

func (act *Account) GetQRCodeSelector() string {
	return ""
}

func (act *Account) SaveCookies([]byte) error {
	return nil
}

func (act *Account) GetCookies() []byte {
	return nil
}

func (act *Account) GetLastURL() string {
	return ""
}

func (act *Account) SetLastURL(url string) {
}

func (act *Account) IsAvailable() bool {
	act.mu.Lock()
	defer act.mu.Unlock()

	return act.curState == stateInitialized
}

func (act *Account) GetMsgChan() chan *message.AsyncMessage {
	return act.msgChan
}

func (act *Account) GetID() int {
	return act.ID
}
