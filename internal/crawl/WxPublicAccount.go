package crawl

type Account1 struct {
}

// TODO plugin
func (a *Account) login() error {
	//cfg, ok := config.AccountConfigMap["WX"]
	//if !ok {
	//	return fmt.Errorf("miss config")
	//}
	//
	//var lastLoggedPath = "https://mp.weixin.qq.com/cgi-bin/home?t=home/index&lang=zh_CN&token=2098303583"
	//if err := util.WxQRCodeLogin("https://mp.weixin.qq.com", "cookies.tmp", "code.png", &lastLoggedPath); err != nil {
	//	return err
	//}

	// get cookie to colly

	return nil
}

func (a *Account) Crawl1(instance *Instance) error {
	return nil
}

//construct login request param
//handle login response
