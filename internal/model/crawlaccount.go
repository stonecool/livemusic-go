package model

import (
	"gorm.io/gorm"
)

type CrawlAccount struct {
	Model

	CrawlType   string // 微信公众号，微博
	AccountId   string
	AccountName string
	Cookies     []byte
}

// AddCrawlAccount Adds a new crawl account
func AddCrawlAccount(data map[string]interface{}) (*CrawlAccount, error) {
	account := CrawlAccount{
		CrawlType: data["crawl_type"].(string),
	}

	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

// GetCrawlAccount Gets a crawl account
func GetCrawlAccount(id int) (*CrawlAccount, error) {
	var account CrawlAccount
	if err := db.Where("id = ?", id).First(&account).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &account, nil
}

// GetCrawlAccountsByType Gets crawl accounts by type
func GetCrawlAccountsByType(crawlType string) ([]*CrawlAccount, error) {
	var crawls []CrawlAccount
	if err := db.Where("deleted_at != 0 AND account_type = ?", crawlType).Find(&crawls).Error; err != nil {
		return nil, err
	}

	ret := make([]*CrawlAccount, len(crawls), len(crawls))
	for i, _ := range crawls {
		ret[i] = &crawls[i]
	}

	return ret, nil
}
