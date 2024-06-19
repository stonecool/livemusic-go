package model

import (
	"gorm.io/gorm"
)

type Crawl struct {
	Model

	CrawlType   string
	AccountId   string
	AccountName string
	Cookies     []byte
}

// AddCrawl Adds a new crawl
func AddCrawl(data map[string]interface{}) (*Crawl, error) {
	crawl := Crawl{
		CrawlType: data["crawl_type"].(string),
	}

	if err := db.Create(&crawl).Error; err != nil {
		return nil, err
	}

	return &crawl, nil
}

// GetCrawlByID Gets a crawl by id
func GetCrawlByID(id int) (*Crawl, error) {
	var crawl Crawl
	if err := db.Where("id = ?", id).First(&crawl).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &crawl, nil
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
