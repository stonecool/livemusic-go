package account

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account/types"
	"github.com/stonecool/livemusic-go/internal/message"
	"github.com/stonecool/livemusic-go/internal/task"
	"time"
)

type account struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	lastURL      string
	cookies      []byte
	instanceID   int
	stateHandler types.StateHandler
	msgChan      chan *message.AsyncMessage
	done         chan struct{}
}

var _ types.Account = (*account)(nil)

func (a *account) Initialize() {
	go a.processTask()
}

func (a *account) processTask() {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case msg := <-a.GetMsgChan():
			a.stateHandler.Transit(a, msg.Cmd)
			if msg.Result != nil {
				msg.Result <- nil
				close(msg.Result)
			}

		case <-ticker.C:
			// 只有在 Ready 状态才尝试获取任务
			if a.stateHandler.GetState() == message.AccountState_AS_Ready {
				_, err := task.DefaultQueue.PopTaskByCategory(a.Category)
				if err != nil {
					continue // 队列为空或没有匹配的任务，继续等待
				}

				// 执行任务
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

func (a *account) GetCategory() string {
	return a.Category
}
