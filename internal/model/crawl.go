package model

import (
	"gorm.io/gorm"
)

type Crawl struct {
	Model

	CrawlType   string // 微信公众号，微博
	AccountId   string
	AccountName string
	Cookies     []byte
}

// AddCrawl Adds a new crawl
func AddCrawl(data map[string]interface{}) (*Crawl, error) {
	c := Crawl{
		CrawlType: data["crawl_type"].(string),
	}

	if err := db.Create(&c).Error; err != nil {
		return nil, err
	}

	return &c, nil
}

// GetCrawlByID Gets a crawl by id
func GetCrawlByID(id int) (*Crawl, error) {
	var c Crawl
	if err := db.Where("id = ?", id).First(&c).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &c, nil
}

// GetCrawlsByType Gets crawls by type
func GetCrawlsByType(crawlType string) ([]*Crawl, error) {
	var crawls []Crawl
	if err := db.Where("deleted_at != 0 AND account_type = ?", crawlType).Find(&crawls).Error; err != nil {
		return nil, err
	}

	ret := make([]*Crawl, len(crawls), len(crawls))
	for i, _ := range crawls {
		ret[i] = &crawls[i]
	}

	return ret, nil
}
