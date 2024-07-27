package model

import (
	"time"
)

type CrawlMsg struct {
	Model

	DataType        string // Livehouse...
	DataId          int
	AccountType     string // 微信公众号，微博
	TargetAccountId string
	Mark            string
	Count           int
	FirstTime       int
	LastTime        int
}

// AddCrawlMsg Adds a crawl message
func AddCrawlMsg(data map[string]interface{}) (*CrawlMsg, error) {
	msg := CrawlMsg{
		DataType:        data["data_type"].(string),
		DataId:          data["data_id"].(int),
		AccountType:     data["account_type"].(string),
		TargetAccountId: data["account_id"].(string),
	}

	if err := db.Create(&msg).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}

// GetCrawlMg Gets a crawl msg
func GetCrawlMg(id int) (*CrawlMsg, error) {
	var msg CrawlMsg
	if err := db.Where("id = ? AND deleted_at = ?", id, 0).First(&msg).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}

// GetCrawlMg Gets a crawl msg
func GetCrawlMsgAll() ([]*CrawlMsg, error) {
	var msgs []*CrawlMsg
	if err := db.Where("deleted_at = ?", 0).Find(&msgs).Error; err != nil {
		return nil, err
	}

	return msgs, nil
}

// DeleteCrawlMsg Deletes a crawl account
func DeleteCrawlMsg(msg *CrawlMsg) error {
	return db.Model(msg).Where("deleted_at != ?", 0).Update("deleted_at", time.Now().Unix()).Error
}

func CrawlMsgExist(dataType string, dataId int, crawlType string) (bool, error) {
	var exists bool
	if err := db.Where("data_type = '?' AND data_id = ? AND crawl_type = '?' AND deleted_at = ?",
		dataType, dataId, crawlType, 0).Find(&exists).Error; err != nil {
		return false, err
	}

	return exists, nil
}

func EditCrawlMsg(id int, data map[string]interface{}) (*CrawlMsg, error) {
	var msg CrawlMsg

	if err := db.Model(&msg).Where("id = ? AND deleted_at = ? ", id, 0).Updates(data).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}
