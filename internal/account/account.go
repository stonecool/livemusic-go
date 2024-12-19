package account

import (
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

			if err = a.stateHandler.HandleStateTransition(msg.Cmd); err != nil {
				a.stateHandler.HandleError(err)
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

func (a *account) Close() {
	close(a.done)
}

func (a *account) GetID() int {
	return a.ID
}

func (a *account) GetName() string {
	return a.Name
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

func (a *account) GetMsgChan() chan *message.AsyncMessage {
	return a.msgChan
}

func (a *account) IsAvailable() bool {
	return a.stateHandler.GetState() == message.AccountState_Ready
}

func (a *account) GetCategory() string {
	return a.Category
}
