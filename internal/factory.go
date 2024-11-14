package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/crawlaccount"
	"log"
)

var crawlCache *cache.Memo

func init() {
	crawlCache = cache.New(getCrawl)
}

func getCrawl(id int) (interface{}, error) {
	account := &crawlaccount.CrawlAccount{ID: id}
	err := account.Get()
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	cfg, ok := config.AccountMap[account.Category]
	if !ok {
		return nil, error(nil)
	}

	var crawl ICrawl
	switch account.Category {
	case "wechat":
		crawl = &WeChatCrawl{
			Crawl: crawl2.Crawl{
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

func startCrawl(crawl crawlaccount.ICrawlAccount) {
	log.Printf("Start account:%d\n", crawl.GetId())

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
	//log.SetOutput(io.Discard)

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

				if err := crawl.Login(); err != nil {
					fmt.Printf("login error:%v", err)
				} else {
					crawl.SetState(CrawlState_Ready)
				}

			case CrawlCmd_Crawl:
				if curState != CrawlState_Ready {
					log.Printf("state not ready")
					continue
				}
				GoCrawl(crawl, crawl.callback)

			default:
				log.Printf("cmd:%v not supportted", clientMessage.message.Cmd)
			}
		}
	}
}

func initialCrawl(account crawlaccount.ICrawlAccount) bool {
	if len(account.GetCookies()) == 0 || len(account.GetLoginURL()) == 0 {
		return false
	}

	err := chromedp.Run(account.GetContext(),
		SetCookies(account),
		chromedp.Navigate(account.GetLastLoginURL()),
		account.CheckLogin(),
	)

	if err != nil {
		log.Printf("%v\n", err)
	}

	return err == nil
}

func GoCrawl(crawl ICrawl, callback Callback) bool {
	err := chromedp.Run(crawl.GetContext(),
		network.Enable(),
		SetCookies(crawl),
		chromedp.Navigate(crawl.GetLastLoginURL()),
		crawl.CheckLogin(),
		crawl.GoCrawl(callback),
	)

	if err != nil {
		log.Printf("%v\n", err)
	}

	return err == nil
}
