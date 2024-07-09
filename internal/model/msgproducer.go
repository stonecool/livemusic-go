package model

import "gorm.io/gorm"

type MsgProducer struct {
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

// AddMsgProducer Adds a crawl message producer
func AddMsgProducer(data map[string]interface{}) (*MsgProducer, error) {
	c := MsgProducer{
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

// GetMsgProducer Gets a crawl coroutine by id
func GetMsgProducer(id int) (*MsgProducer, error) {
	var m MsgProducer
	if err := db.Where("id = ?", id).First(&m).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &m, nil
}

// MsgProducerExists Check coroutine exists
func MsgProducerExists(dataType string, dataId int, crawlType string) bool {
	var exists bool
	if err := db.Where("data_type = '?' AND data_id = ? AND crawl_type = '?'", dataType, dataId, crawlType).Find(&exists).Error; err != nil {
		return false
	}

	return exists
}
