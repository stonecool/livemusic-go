package factory

import (
	"fmt"
)

// GetFactory 根据类别返回相应的工厂
func GetFactory(category string) (ICrawlFactory, error) {
	switch category {
	case "wechat":
		return &WeChatFactory{}, nil
	// 可以添加其他类型的工厂
	default:
		return nil, fmt.Errorf("unsupported category: %s", category)
	}
}

//func getCrawl(id int) (interface{}, error) {
//	account := &CrawlAccount{ID: id}
//	err := account.Get()
//	if err != nil {
//		return nil, err
//	}
//
//	factory, err := GetFactory(account.Category)
//	if err != nil {
//		return nil, err
//	}
//
//	cfg, ok := config.AccountMap[account.Category]
//	if !ok {
//		return nil, fmt.Errorf("config not found for category: %s", account.Category)
//	}
//
//	crawl := factory.CreateCrawl(&cfg)
//	go startCrawl(crawl)
//	return crawl, nil
//}
