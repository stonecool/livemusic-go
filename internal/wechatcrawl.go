package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"net/url"
	"strings"
	"time"
)

type WeChatCrawl struct {
	Crawl
}

func (crawl *WeChatCrawl) GetQRCodeSelector() string {
	return "#header > div.banner > div > div > div.login__type__container.login__type__container__scan > img"
}

func (crawl *WeChatCrawl) SetLastLoginURL(url string) {
	crawl.Account.lastLoginURL = url
}

func (crawl *WeChatCrawl) GetLastLoginURL() string {
	return crawl.Account.lastLoginURL
}

func (crawl *WeChatCrawl) WaitLogin() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		return chromedp.WaitVisible(`#app > div.main_bd_new`, chromedp.ByID).Do(ctx)
	}
}

func (crawl *WeChatCrawl) CheckLogin() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		var msg string
		if err = chromedp.Evaluate(`document.querySelector('#body > div > div > div > div > div > div').className`, &msg).Do(ctx); err != nil {
			return
		}

		if msg == "page_error_msg" {
			return fmt.Errorf("login error")
		}

		return
	}
}

func (crawl *WeChatCrawl) GoCrawl() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		var currentURL string
		err = chromedp.Location(&currentURL).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}

		// 解析 URL 并获取参数键值对
		parsedURL, err := url.Parse(currentURL)
		if err != nil {
			log.Fatal(err)
		}

		var token string
		queryParams := parsedURL.Query()
		for key, values := range queryParams {
			for _, value := range values {
				if key == "token" {
					token = value
					break
				}
			}

			if key == "token" {
				break
			}
		}

		searchFakeIdURLTemplate := "https://mp.weixin.qq.com/cgi-bin/searchbiz?action=search_biz&begin=%d&count=%d&query=%s&token=%s&lang=zh_CN&f=json&ajax=1"
		searchFakeIdURL := fmt.Sprintf(searchFakeIdURLTemplate, 0, 10, "1701MusicPark", token)

		requestExpression := `(async function() {
			const response = await fetch('%s');
			const data = await response.json();
			return data;
       })();
`
		type BaseResp struct {
			Ret    int    `json:"ret"`
			ErrMsg string `json:"err_msg"`
		}

		type SearchNameResp struct {
			BaseResp BaseResp `json:"base_resp"`
			List     []struct {
				FakeId   string `json:"fakeid"`
				Nickname string `json:"nickname"`
			} `json:"list"`
		}

		var searchNameResp SearchNameResp
		if err = chromedp.Evaluate(fmt.Sprintf(requestExpression, searchFakeIdURL), &searchNameResp, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}).Do(ctx); err != nil {
			fmt.Printf("%v\n", err)
			return err
		}

		if searchNameResp.BaseResp.Ret != 0 || searchNameResp.BaseResp.ErrMsg != "ok" {
			return fmt.Errorf("base resp error")
		}

		fakeId := ""
		for _, item := range searchNameResp.List {
			if item.Nickname == crawl.Account.AccountName {
				fakeId = item.FakeId
				break
			}
		}

		if fakeId == "" {
			return fmt.Errorf("not found fakeid")
		}

		getURLTemplate := "https://mp.weixin.qq.com/cgi-bin/appmsgpublish?sub=list&search_field=null&begin=%d&count=%d&query=&fakeid=%s&type=101_1&free_publish_type=1&sub_action=list_ex&token=%s&lang=zh_CN&f=json&ajax=1"
		getURL := fmt.Sprintf(getURLTemplate, 0, 5, fakeId, token)

		type GetListResp struct {
			BaseResp    BaseResp `json:"base_resp"`
			PublishPage string   `json:"publish_page"`
		}
		var getListResp GetListResp

		if err = chromedp.Evaluate(fmt.Sprintf(requestExpression, getURL), &getListResp, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}).Do(ctx); err != nil {
			return err
		}

		if getListResp.BaseResp.Ret != 0 || getListResp.BaseResp.ErrMsg != "ok" {
			return fmt.Errorf("base resp error")
		}

		type PublishPage struct {
			PublishList []struct {
				PublishInfo string `json:"publish_info"`
			} `json:"publish_list"`
		}

		type PublishInfo struct {
			AppMsgEx []struct {
				Title string `json:"title"`
				Link  string `json:"link"`
			} `json:"appmsgex"`
		}

		var publishPage PublishPage
		err = json.Unmarshal([]byte(getListResp.PublishPage), &publishPage)
		if err != nil {
			return err
		}

		chromedp.ListenTarget(ctx, func(ev interface{}) {
			switch ev := ev.(type) {
			case *network.EventResponseReceived:
				// 检查请求的 URL
				if strings.HasPrefix(ev.Response.URL, "https://mp.weixin.qq.com/mp/appmsg_weapp") {
					go func(requestID network.RequestID) {
						responseBody, err := network.GetResponseBody(requestID).Do(ctx)
						if err != nil {
							log.Printf("Failed to get response body: %v\n", err)
							return
						}

						type DetailResp struct {
							AppMsgCompatURL string `json:"appmsg_compact_url"`
						}

						var detailResp DetailResp
						if err := json.Unmarshal(responseBody, &detailResp); err != nil {
							log.Printf("Failed to unmarshal JSON response: %v\n", err)
							return
						}

						fmt.Printf("Compact: %s\n", detailResp.AppMsgCompatURL)
					}(ev.RequestID)
				}
			}
		})

		for _, item := range publishPage.PublishList {
			var publishInfo PublishInfo

			err = json.Unmarshal([]byte(item.PublishInfo), &publishInfo)
			if err != nil {
				continue
			}

			for _, msg := range publishInfo.AppMsgEx {
				fmt.Printf("Title: %s Link: %s\n", msg.Title, msg.Link)

				err = chromedp.Navigate(msg.Link).Do(ctx)
				if err != nil {
					log.Fatal(err)
				}

				err = chromedp.Sleep(5 * time.Second).Do(ctx) // 等待一段时间以确保请求完成
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		return
	}
}
