package model

import "gorm.io/gorm"

type CrawlMsg struct {
	Model

	DataType  string // Livehouse...
	DataId    int
	CrawlType string // 微信公众号，微博
	AccountId string
	Mark      string
	Count     int
	FirstTime int
	LastTime  int
}

// AddCrawlMsg Adds a crawl message producer
func AddCrawlMsg(data map[string]interface{}) (*CrawlMsg, error) {
	c := CrawlMsg{
		DataType:  data["data_type"].(string),
		DataId:    data["data_id"].(int),
		CrawlType: data["crawl_type"].(string),
		AccountId: data["account_id"].(string),
	}

	if err := db.Create(&c).Error; err != nil {
		return nil, err
	}

	return &c, nil
}

// GetCrawlMsg Gets a crawl coroutine by id
func GetCrawlMsg(id int) (*CrawlMsg, error) {
	var m CrawlMsg
	if err := db.Where("id = ?", id).First(&m).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &m, nil
}

// CrawlMsgExists Check coroutine exists
func CrawlMsgExists(dataType string, dataId int, crawlType string) bool {
	var exists bool
	if err := db.Where("data_type = '?' AND data_id = ? AND crawl_type = '?'", dataType, dataId, crawlType).Find(&exists).Error; err != nil {
		return false
	}

	return exists
}
