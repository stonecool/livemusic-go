package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/config"
	"io"
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

func startCrawl(crawl ICrawl) {
	log.Printf("Start crawl:%d\n", crawl.GetId())

	log.SetOutput(io.Discard)
	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),

		append(
			chromedp.DefaultExecAllocatorOptions[:],
			//chromedp.NoDefaultBrowserCheck,
			chromedp.Flag("headless", false),
			//chromedp.Flag("hide-scrollbars", false),
			//chromedp.Flag("mute-audio", false),
			//chromedp.Flag("ignore-certificate-errors", true),
			//chromedp.Flag("disable-web-security", true),
			//chromedp.Flag("disable-gpu", false),
			//chromedp.NoFirstRun,
			//chromedp.Flag("enable-automation", false),
			//chromedp.Flag("disable-extensions", false),
		)...,
	)

	defer cancel()

	// create a timeout
	//ctx, cancel := context.WithTimeout(ctx, 150*time.Second)
	//defer cancel()

	// create chrome instance
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	crawl.SetContext(ctx)

	for {
		select {
		case clientMessage := <-crawl.GetChan():
			curState := crawl.GetState()
			switch clientMessage.message.Cmd {
			case CrawlCmd_Initial:
				if curState != CrawlState_Uninitialized {
					continue
				}

				if initialCrawl(crawl) {
					crawl.SetState(CrawlState_Ready)
				} else {
					crawl.SetState(CrawlState_NotLogged)
				}

			case CrawlCmd_Login:
				if curState != CrawlState_NotLogged {
					log.Printf("state not ready")
					continue
				}

				if err := loginCrawl(crawl); err != nil {
					fmt.Printf("login error:%v", err)
				} else {
					crawl.SetState(CrawlState_Ready)
				}

			case CrawlCmd_Crawl:
				if curState != CrawlState_Ready {
					log.Printf("state not ready")
					continue
				}
				goCrawl(crawl)

			default:
				log.Printf("cmd:%v not supportted", clientMessage.message.Cmd)
			}
		}
	}
}

func initialCrawl(crawl ICrawl) bool {
	if len(crawl.GetCookies()) == 0 || len(crawl.GetLastLoginURL()) == 0 {
		return false
	}

	return checkLogin(crawl)
}

func loginCrawl(crawl ICrawl) error {
	return QRCodeLogin(crawl)
}

func goCrawl(crawl ICrawl) {
	GoCrawl1(crawl)
}
