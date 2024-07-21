package model

import "gorm.io/gorm"

type CrawlMsg struct {
	Model

	DataType    string // Livehouse...
	DataId      int
	AccountType string // 微信公众号，微博
	AccountId   string
	Mark        string
	Count       int
	FirstTime   int
	LastTime    int
}

// AddCrawlMsg Adds a crawl message
func AddCrawlMsg(data map[string]interface{}) (*CrawlMsg, error) {
	msg := CrawlMsg{
		DataType:    data["data_type"].(string),
		DataId:      data["data_id"].(int),
		AccountType: data["account_type"].(string),
		AccountId:   data["account_id"].(string),
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

// GetCrawlMg Gets a crawl msg
func GetCrawlMsgAll() ([]*CrawlMsg, error) {
	var msgs []*CrawlMsg
	if err := db.Where("deleted_at != ?", 0).Find(&msgs).Error; err != nil {
		return nil, err
	}

	return msgs, nil
}

// DeleteCrawlMsg Deletes a crawl account
func DeleteCrawlMsg(msg *CrawlMsg) error {
	return db.Delete(msg).Error
}

// CrawlMsgExists Check coroutine exists
func CrawlMsgExists(dataType string, dataId int, crawlType string) bool {
	var exists bool
	if err := db.Where("data_type = '?' AND data_id = ? AND crawl_type = '?'", dataType, dataId, crawlType).Find(&exists).Error; err != nil {
		return false
	}

	return exists
}

func EditCrawlMsg(id int, data map[string]interface{}) (*CrawlMsg, error) {
	var msg CrawlMsg

	if err := db.Model(&msg).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}
