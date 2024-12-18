package account

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account/state"
	"github.com/stonecool/livemusic-go/internal/message"
)

type account struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	lastURL      string
	cookies      []byte
	instanceID   int
	stateHandler state.Handler
	msgChan      chan *message.AsyncMessage
	done         chan struct{}
}

func (a *account) processTask() {
	for {
		select {
		case msg := <-a.msgChan:
			var err error

			if err = a.handleCommand(msg); err != nil {
				a.stateHandler.HandleError(err)
			} else {
				if err = a.stateHandler.HandleStateTransition(msg.Cmd); err != nil {
					a.stateHandler.HandleError(err)
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

func (a *account) handleCommand(msg *message.AsyncMessage) error {
	currentState := a.stateHandler.GetState()

	switch currentState {
	case message.AccountState_New:
		return fmt.Errorf("account not initialized")

	case message.AccountState_NotLoggedIn:
		if msg.Cmd != message.AccountCmd_Login {
			return fmt.Errorf("invalid command:%v for not logged in state", msg.Cmd)
		}
		return a.handleLogin()

	case message.AccountState_Ready:
		switch msg.Cmd {
		case message.AccountCmd_Crawl:
			return a.handleCrawl(msg.Data)
		case message.AccountCmd_Login:
			return a.handleLogin()
		default:
			return fmt.Errorf("invalid command:%v for ready state", msg.Cmd)
		}

	case message.AccountState_Running:
		return fmt.Errorf("account is busy")

	case message.AccountState_Expired:
		if msg.Cmd != message.AccountCmd_Login {
			return fmt.Errorf("invalid command:%v for expired state", msg.Cmd)
		}
		return a.handleLogin()

	default:
		return fmt.Errorf("invalid state: %v", currentState)
	}
}

func (a *account) Close() {
	close(a.done)
}

func (a *account) handleLogin() error {
	// 处理登录任务的具体逻辑
	return nil
}

func (a *account) handleCrawl(payload interface{}) error {
	// 处理爬取任务的具体逻辑
	return nil
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

func (a *account) IsAvailable() bool {
	return a.stateHandler.GetState() == message.AccountState_Ready
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
