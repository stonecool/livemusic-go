package account

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account/types"
	"github.com/stonecool/livemusic-go/internal/message"
	"go.uber.org/zap"
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

// 确保 account 类型实现了 interfaces.Account 接口
var _ types.Account = (*account)(nil)

func (a *account) processTask() {
	for {
		select {
		case msg := <-a.msgChan:
			a.stateHandler.Transit(a, msg.Cmd)

			err := a.executeCommand(msg.Cmd)
			if err != nil {
				internal.Logger.Error("Account command execution failed",
					zap.Int("account_id", a.ID),
					zap.String("command", msg.Cmd.String()),
					zap.Error(err),
				)
				a.stateHandler.HandleError(a, err)
				if msg.Result != nil {
					msg.Result <- err
					close(msg.Result)
				}
				continue
			}

			internal.Logger.Info("Account command completed",
				zap.Int("account_id", a.ID),
				zap.String("command", msg.Cmd.String()),
			)

			if msg.Result != nil {
				msg.Result <- nil
				close(msg.Result)
			}

		case <-a.done:
			return
		}
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
	// 启动一个新��� goroutine 来执行爬虫任务
	go func() {
		// 执行具体的爬虫逻辑
		//err := a.executeCrawlTask()

		// 任务完成后发送 CrawlAck 消息
		result := make(chan error, 1)
		//a.msgChan <- &message.AsyncMessage{
		//	Cmd:    message.AccountCmd_CrawlAck,
		//	Result: result,
		//}

		// 等待状态转换完成
		<-result
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
