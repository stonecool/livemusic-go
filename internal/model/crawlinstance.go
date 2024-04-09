package model

import (
	"gorm.io/gorm"
)

type CrawlInstance struct {
	Model

	Name        string `json:"name"`
	AccountId   int    `json:"account_id"`
	Headers     string `json:"headers"`
	QueryParams string `json:"query_params"`
	FormData    string `json:"form_data"`
	State       uint8  `json:"state"`
	Count       uint   `json:"count"`
	FirstTime   uint   `json:"first_time"`
	LastTime    uint   `json:"last_time"`
}

// AddCrawlInstance Adds a new crawl
func AddCrawlInstance(data map[string]interface{}) (int, error) {
	crawl := CrawlInstance{
		Name:      data["name"].(string),
		AccountId: data["account_id"].(int),
		Count:     0,
		FirstTime: 0,
		LastTime:  0,
		State:     0,
	}

	if err := db.Create(&crawl).Error; err != nil {
		return 0, err
	}

	return crawl.ID, nil
}

// GetCrawlInstance Gets a crawl by id
func GetCrawlInstance(id int) (*CrawlInstance, error) {
	var crawl CrawlInstance
	if err := db.Where("id = ?", id).First(&crawl).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &crawl, nil
}

// GetCrawlInstances Get all crawls
func GetCrawlInstances() ([]*CrawlInstance, error) {
	var crawls []CrawlInstance
	if err := db.Where("id <= ?", 10).Find(&crawls).Error; err != nil {
		return nil, err
	}

	ret := make([]*CrawlInstance, len(crawls), len(crawls))
	for i, s := range crawls {
		ret[i] = &s
	}
	return ret, nil
}
