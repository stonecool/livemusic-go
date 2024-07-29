package internal

import (
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/config"
	"log"
)

var crawlCache *cache.Memo

func init() {
	crawlCache = cache.New(getCrawl)
}

func getCrawl(id int) (interface{}, error) {
	account := &CrawlAccount{ID: id}
	err := account.Get()
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	cfg, ok := config.AccountMap[account.AccountType]
	if !ok {
		return nil, error(nil)
	}

	var crawl ICrawl
	switch account.AccountType {
	case "WeChat":
		crawl = &WeChatCrawl{
			Crawl: Crawl{
				Account: account,
				state:   CrawlState_Uninitialized,
				config:  &cfg,
				ch:      make(chan *ClientMessage),
			},
		}
	}

	go startCrawl(crawl)
	return crawl, nil
}

func GetCrawl(id int) (ICrawl, error) {
	crawl, err := crawlCache.Get(id)
	if err != nil {
		return nil, err
	} else {
		return crawl.(ICrawl), nil
	}
}

func startCrawl(c ICrawl) {
	log.Printf("Start c:%d\n", c.GetId())

	for {
		select {
		case clientMessage := <-c.GetChan():
			curState := c.GetState()
			switch clientMessage.message.Cmd {
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
				log.Printf("cmd:%v not supportted", clientMessage.message.Cmd)
			}
		}
	}
}
