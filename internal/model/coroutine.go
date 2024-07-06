package model

import "gorm.io/gorm"

type CrawlCoroutine struct {
	Model

	DataType  string // Livehouse...
	DataId    int
	CrawlType string // 微信公众号，微博
	AccountId int
	CrawlMark string
	Count     int
	FirstTime int
	LastTime  int
}

// AddCrawlCoroutine Adds a crawl coroutine
func AddCrawlCoroutine(data map[string]interface{}) (*CrawlCoroutine, error) {
	c := CrawlCoroutine{}

	if err := db.Create(&c).Error; err != nil {
		return nil, err
	}

	return &c, nil
}

// GetCrawlCoroutine Gets a crawl coroutine by id
func GetCrawlCoroutine(id int) (*CrawlCoroutine, error) {
	var c CrawlCoroutine
	if err := db.Where("id = ?", id).First(&c).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &c, nil
}
