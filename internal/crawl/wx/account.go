package wx

import (
	"fmt"
	"github.com/stonecool/1701livehouse-server/internal/config"
	"github.com/stonecool/1701livehouse-server/internal/crawl"
	"github.com/stonecool/1701livehouse-server/internal/util"
)

type Account struct {
	crawl.Account
}

// TODO plugin
func (a *Account) login() error {
	cfg, ok := config.AccountTemplateMap["WX"]
	if !ok {
		return fmt.Errorf("miss config")
	}

	var lastLoggedPath = "https://mp.weixin.qq.com/cgi-bin/home?t=home/index&lang=zh_CN&token=2098303583"
	if err := util.WxQRCodeLogin("https://mp.weixin.qq.com", "cookies.tmp", "code.png", &lastLoggedPath); err != nil {
		return err
	}

	// get cookie to colly

	return nil
}

func (a *Account) Crawl(instance *crawl.Instance) error {
	return nil
}

//construct login request param
//handle login response
