package model

import (
	"gorm.io/gorm"
	"time"
)

type CrawlRoutine struct {
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

// AddCrawlRoutine Adds a crawl message
func AddCrawlRoutine(data map[string]interface{}) (*CrawlRoutine, error) {
	routine := CrawlRoutine{
		DataType:        data["data_type"].(string),
		DataId:          data["data_id"].(int),
		AccountType:     data["account_type"].(string),
		TargetAccountId: data["target_account_id"].(string),
	}

	if err := db.Create(&routine).Error; err != nil {
		return nil, err
	}

	return &routine, nil
}

// GetCrawlRoutine Gets a crawl msg
func GetCrawlRoutine(id int) (*CrawlRoutine, error) {
	var routine CrawlRoutine
	if err := db.Where("id = ? AND deleted_at = ?", id, 0).First(&routine).Error; err != nil {
		return nil, err
	}

	return &routine, nil
}

func GetCrawlRoutineAll() ([]*CrawlRoutine, error) {
	var s []*CrawlRoutine
	if err := db.Where("deleted_at = ?", 0).Find(&s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteCrawlRoutine Deletes a crawl account
func DeleteCrawlRoutine(msg *CrawlRoutine) error {
	return db.Model(msg).Where("deleted_at = ?", 0).Update("deleted_at", time.Now().Unix()).Error
}

func ExistCrawlRoutine(dataType string, dataId int, crawlType string) (bool, error) {
	var routine CrawlRoutine
	err := db.Select("id").Where("data_type = ? AND data_id = ? AND account_type = ? AND deleted_at = ?",
		dataType, dataId, crawlType, 0).First(&routine).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return routine.ID > 0, nil
}

func EditCrawlRoutine(id int, data map[string]interface{}) (*CrawlRoutine, error) {
	var routine CrawlRoutine

	if err := db.Model(&routine).Where("id = ? AND deleted_at = ? ", id, 0).Updates(data).Error; err != nil {
		return nil, err
	}

	return &routine, nil
}
