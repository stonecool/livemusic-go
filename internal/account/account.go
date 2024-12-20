package account

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account/types"
	"github.com/stonecool/livemusic-go/internal/message"
	"go.uber.org/zap"
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

func (a *account) processTask() {
	ticker := time.NewTicker(time.Second)
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
			if a.stateHandler.GetState() == message.AccountState_Ready {
				task, err := message.DefaultQueue.PopTaskByCategory(a.Category)
				if err != nil {
					continue // 队列为空或没有匹配的任务，继续等待
				}

				// 执行任务
				if err := a.executeCommand(task.Message.Cmd); err != nil {
					internal.Logger.Error("Failed to execute task",
						zap.Int("account_id", a.ID),
						zap.String("category", a.Category),
						zap.Error(err))
					a.stateHandler.HandleError(a, err)
					if task.Message.Result != nil {
						task.Message.Result <- err
						close(task.Message.Result)
					}
				}
			}

		case <-a.done:
			return
		}
	}
}

func (a *account) canExecuteCommand(cmd message.AccountCmd) bool {
	currentState := a.stateHandler.GetState()
	switch cmd {
	case message.AccountCmd_Crawl:
		return currentState == message.AccountState_Ready
	case message.AccountCmd_CrawlAck:
		return currentState == message.AccountState_Running
	case message.AccountCmd_Login:
		return currentState == message.AccountState_NotLoggedIn
	default:
		return true
	}
}

func (a *account) executeCommand(cmd message.AccountCmd) error {
	switch cmd {
	case message.AccountCmd_Login:
		return a.Login()
	case message.AccountCmd_Crawl:
		return a.doCrawl()
	case message.AccountCmd_CrawlAck:
		return nil
	default:
		return fmt.Errorf("unknown command: %v", cmd)
	}
}

func (a *account) doCrawl() error {
	go func() {
		err := a.executeCrawlTask()

		// 爬虫完成后，发送状态转换命令到控制通道
		result := make(chan error, 1)
		ackMsg := message.NewAsyncMessage(message.AccountCmd_CrawlAck, result)

		select {
		case a.controlChan <- ackMsg:
			<-result
		case <-a.done:
			return
		case <-time.After(5 * time.Second):
			internal.Logger.Error("Failed to send CrawlAck command",
				zap.Int("account_id", a.ID))
		}
	}()
	return nil
}

func (a *account) executeCrawlTask() error {
	// TODO: 实现具体的爬虫逻辑
	return nil
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
