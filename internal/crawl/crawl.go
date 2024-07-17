package crawl

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/config"
	"log"
)

type Crawl struct {
	Account *internal.CrawlAccount

	state  internal.CrawlState
	config *config.Account
	ch     chan *internal.Message
}

func (crawl *Crawl) GetId() string {
	return crawl.Account.AccountId
}

func (crawl *Crawl) GetName() string {
	return crawl.Account.AccountName
}

func (crawl *Crawl) GetState() internal.CrawlState {
	return internal.CrawlState_Ready
}

func (crawl *Crawl) SetState(state internal.CrawlState) {
	//crawl.State = state.
}

func (crawl *Crawl) Login() (bool, error) {
	return false, nil
}

func (crawl *Crawl) CheckLogin() chromedp.ActionFunc {
	return nil
}
func (crawl *Crawl) GetLoginURL() string {
	return crawl.config.LoginURL
}

func (crawl *Crawl) GetQRCode(data []byte) {

}

func (crawl *Crawl) GetQRCodeSelector() string {
	return ""
}

func (crawl *Crawl) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (crawl *Crawl) SaveCookies([]byte) {
}

func (crawl *Crawl) GetCookies() []byte {
	return nil
}

func (crawl *Crawl) GetChan() chan *internal.Message {
	return nil
}

func (crawl *Crawl) Start() {
	log.Printf("Start crawl:%d\n", crawl.GetId())

	for {
		select {
		case msg := <-crawl.GetChan():
			curState := crawl.GetState()

			switch msg.Cmd {
			case internal.CrawlCmd_Initial:
				if curState != internal.CrawlState_Uninitialized {
					continue
				}

				ret, err := crawl.Login()
				if err != nil {
					log.Printf("error:%s", err)
					continue
				}

				if ret {
					crawl.SetState(internal.CrawlState_NotLogged)
				}

			case internal.CrawlCmd_Login:
				if curState != internal.CrawlState_NotLogged {
					log.Printf("state not ready")
					continue
				}

				crawl.SetState(internal.CrawlState_Ready)

			case internal.CrawlCmd_Crawl:

			default:
				log.Printf("cmd:%v not supportted", msg.Cmd)
			}
		}
	}
}
