package crawl

type WxPublicAccountCrawl struct {
	Crawl
}

func (c *WxPublicAccountCrawl) CheckLogin() (bool, error) {
	return false, nil
}

//func (c *WxPublicAccountCrawl) CheckLogin(URL *string) chromedp.ActionFunc {
//	return func(ctx context.Context) (err error) {
//		chromedp.Navigate(*URL)
//
//		var url string
//		if err = chromedp.Evaluate(`document.querySelector('#body > div > div > div > div > div > div').className`, &url).Do(ctx); err != nil {
//			return
//		}
//
//		if url == "page_error_msg" {
//			log.Println("no login.")
//			return
//
//		}
//
//		if err = chromedp.Evaluate(`window.location.href`, URL).Do(ctx); err != nil {
//			return
//		}
//
//		log.Println("login succeed.")
//		chromedp.Stop()
//
//		return
//	}
//}

func (c *WxPublicAccountCrawl) Login() (bool, error) {
	return false, nil
}
