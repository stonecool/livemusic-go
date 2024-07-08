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
	c := CrawlCoroutine{
		DataType:  data["data_type"].(string),
		DataId:    data["data_id"].(int),
		CrawlType: data["crawl_type"].(string),
		AccountId: data["account_id"].(int),
	}

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

// CrawlCoroutineExists Check coroutine exists
func CrawlCoroutineExists(dataType string, dataId int, crawlType string) bool {
	var exists bool
	if err := db.Where("data_type = '?' AND data_id = ? AND crawl_type = '?'", dataType, dataId, crawlType).Find(&exists).Error; err != nil {
		return false
	}

	return exists
}

// CrawlCoroutineExists Check coroutine exists
func CrawlCoroutineExistsById(id int) bool {
	var exists bool
	if err := db.Where("id = '?'", id).Find(&exists).Error; err != nil {
		return false
	}

	return exists
}
