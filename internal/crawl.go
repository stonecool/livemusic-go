package internal

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/config"
	"log"
)

type Crawl struct {
	Account *CrawlAccount

	state  CrawlState
	config *config.Account
	ch     chan *Message
}

func (c *Crawl) GetId() int {
	return c.Account.ID
}

func (c *Crawl) GetName() string {
	return c.Account.AccountName
}

func (c *Crawl) GetState() CrawlState {
	return CrawlState_Ready
}

func (c *Crawl) SetState(state CrawlState) {
	//c.State = state.
}

func (c *Crawl) Login() (bool, error) {
	return false, nil
}

func (c *Crawl) CheckLogin() chromedp.ActionFunc {
	return nil
}
func (c *Crawl) GetLoginURL() string {
	return c.config.LoginURL
}

func (c *Crawl) GetQRCode(data []byte) {

}

func (c *Crawl) GetQRCodeSelector() string {
	return ""
}

func (c *Crawl) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (c *Crawl) SaveCookies([]byte) {
}

func (c *Crawl) GetCookies() []byte {
	return nil
}

func (c *Crawl) GetChan() chan *Message {
	return c.ch
}

func (c *Crawl) Start() {
	log.Printf("Start c:%d\n", c.GetId())

	for {
		select {
		case msg := <-c.GetChan():
			curState := c.GetState()

			switch msg.Cmd {
			case CrawlCmd_Initial:
				if curState != CrawlState_Uninitialized {
					continue
				}

				ret, err := c.Login()
				if err != nil {
					log.Printf("error:%s", err)
					continue
				}

				if ret {
					c.SetState(CrawlState_NotLogged)
				}

			case CrawlCmd_Login:
				if curState != CrawlState_NotLogged {
					log.Printf("state not ready")
					continue
				}

				c.SetState(CrawlState_Ready)

			case CrawlCmd_Crawl:

			default:
				log.Printf("cmd:%v not supportted", msg.Cmd)
			}
		}
	}
}
