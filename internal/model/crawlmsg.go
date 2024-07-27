package model

import (
	"gorm.io/gorm"
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
		TargetAccountId: data["target_account_id"].(string),
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
	return db.Model(msg).Where("deleted_at = ?", 0).Update("deleted_at", time.Now().Unix()).Error
}

func ExistCrawlMsg(dataType string, dataId int, crawlType string) (bool, error) {
	var msg CrawlMsg
	err := db.Select("id").Where("data_type = ? AND data_id = ? AND account_type = ? AND deleted_at = ?",
		dataType, dataId, crawlType, 0).First(&msg).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return msg.ID > 0, nil
}

func EditCrawlMsg(id int, data map[string]interface{}) (*CrawlMsg, error) {
	var msg CrawlMsg

	if err := db.Model(&msg).Where("id = ? AND deleted_at = ? ", id, 0).Updates(data).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}
