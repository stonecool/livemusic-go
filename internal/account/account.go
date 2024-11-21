package account

import (
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/message"
)

type Account struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	AccountName string `json:"account_name"`
	lastURL     string
	cookies     []byte
	InstanceID  int
	State       state
	mu          sync.RWMutex
	msgChan     chan *message.AsyncMessage
	done        chan struct{}
}

func (acc *Account) Init() {
	go acc.processTask()
	acc.State = stateInitialized
}

func (acc *Account) processTask() {
	for {
		select {
		case msg := <-acc.msgChan:
			currentState := acc.GetState()
			var err error

			switch currentState {
			case stateNew:
				msg.Error = fmt.Errorf("invalid command for new state: %v", msg.Cmd)
				continue
				err = acc.handleInit()

			case stateInitialized:
				// 初始化完成状态只接受登录命令
				if msg.Cmd != message.CrawlCmd_Login {
					msg.Error = fmt.Errorf("invalid command for initialized state: %v", msg.Cmd)
					continue
				}
				err = acc.handleLogin()

			case stateNotLoggedIn:
				// 未登录状态只接受登录命令
				if msg.Cmd != message.CrawlCmd_Login {
					msg.Error = fmt.Errorf("invalid command for not logged in state: %v", msg.Cmd)
					continue
				}
				err = acc.handleLogin()

			case stateReady:
				// 就绪状态可以接受爬取或重新登录命令
				switch msg.Cmd {
				case message.CrawlCmd_Crawl:
					err = acc.handleCrawl(msg.Data)
				case message.CrawlCmd_Login:
					err = acc.handleLogin()
				default:
					msg.Error = fmt.Errorf("invalid command for ready state: %v", msg.Cmd)
					continue
				}

			case stateRunning:
				// 运行状态只能等待当前任务完成
				msg.Error = fmt.Errorf("account is busy")
				continue

			case stateTerminated:
				// 终止状态不接受任何命令
				msg.Error = fmt.Errorf("account is terminated")
				continue
			}

			// 根据操作结果更新状态
			if err != nil {
				msg.Error = err
				// 错误处理可能导致状态变化
				switch currentState {
				case stateInitialized:
					acc.SetState(stateNew)
				case stateNotLoggedIn, stateReady:
					acc.SetState(stateNotLoggedIn)
				case stateRunning:
					acc.SetState(stateReady)
				}
			} else {
				// 成功处理后的状态转换
				switch currentState {
				case stateNew:
					acc.SetState(stateInitialized)
				case stateInitialized, stateNotLoggedIn:
					acc.SetState(stateReady)
				case stateReady:
					if msg.Cmd == message.CrawlCmd_Crawl {
						acc.SetState(stateRunning)
					}
				case stateRunning:
					acc.SetState(stateReady)
				}
			}

			// 设置响应结果
			if msg.Error == nil {
				msg.Data = []byte("success")
			}

		case <-acc.done:
			return
		}
	}
}

func (acc *Account) Close() {
	close(acc.done)
}

func (acc *Account) handleLogin() interface{} {
	// 处理登录任务的具体逻辑
	return nil
}

func (acc *Account) handleCrawl(payload interface{}) interface{} {
	// 处理爬取任务的具体逻辑
	return payload
}

func (acc *Account) Get() error {
	//if crawlcaount, err := database.Getcaount(acc.ID); err != nil {
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

func (acc *Account) GetName() string {
	return acc.AccountName
}

func (acc *Account) GetState() state {
	acc.mu.RLock()
	defer acc.mu.RUnlock()

	return acc.State
}

func (acc *Account) SetState(s state) {
	acc.mu.Lock()
	defer acc.mu.Unlock()

	acc.State = s
}

func (acc *Account) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (acc *Account) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (acc *Account) GetLoginURL() string {
	return ""
}

func (acc *Account) Login() error {
	return nil
}

func (acc *Account) GetQRCode([]byte) {
}

func (acc *Account) GetQRCodeSelector() string {
	return ""
}

func (acc *Account) SaveCookies([]byte) error {
	return nil
}

func (acc *Account) GetCookies() []byte {
	return nil
}

func (acc *Account) GetLastURL() string {
	return ""
}

func (acc *Account) SetLastURL(url string) {
}

func (acc *Account) IsAvailable() bool {
	acc.mu.Lock()
	defer acc.mu.Unlock()

	return acc.State == stateInitialized
}

func (acc *Account) GetMsgChan() chan *message.AsyncMessage {
	return acc.msgChan
}

func (acc *Account) GetID() int {
	return acc.ID
}
