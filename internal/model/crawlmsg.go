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

// AddCrawlMsg Adds a crawl message
func AddCrawlMsg(data map[string]interface{}) (*CrawlMsg, error) {
	msg := CrawlMsg{
		DataType:  data["data_type"].(string),
		DataId:    data["data_id"].(int),
		CrawlType: data["crawl_type"].(string),
		AccountId: data["account_id"].(string),
	}

	if err := db.Create(&msg).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}

// GetCrawlMg Gets a crawl msg
func GetCrawlMg(id int) (*CrawlMsg, error) {
	var msg CrawlMsg
	if err := db.Where("id = ?", id).First(&msg).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &msg, nil
}

// CrawlMsgExists Check coroutine exists
func CrawlMsgExists(dataType string, dataId int, crawlType string) bool {
	var exists bool
	if err := db.Where("data_type = '?' AND data_id = ? AND crawl_type = '?'", dataType, dataId, crawlType).Find(&exists).Error; err != nil {
		return false
	}

	return exists
}
