package model

import (
	"time"

	"gorm.io/gorm"
)

type CrawlTask struct {
	Model

	DataType        string // Livehouse...
	DataId          int
	AccountType     string // 微信公众号，微博
	TargetAccountId string
	Mark            string
	Count           int
	FirstTime       int
	LastTime        int
	CronSpec        string // cron表达式,如 "0 0 9 * * *"
}

func AddCrawlTask(data map[string]interface{}) (*CrawlTask, error) {
	task := CrawlTask{
		DataType:        data["data_type"].(string),
		DataId:          data["data_id"].(int),
		AccountType:     data["account_type"].(string),
		TargetAccountId: data["target_account_id"].(string),
		CronSpec:        data["cron_spec"].(string),
	}

	if err := DB.Create(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func GetCrawlTask(id int) (*CrawlTask, error) {
	var task CrawlTask
	if err := DB.Where("id = ? AND deleted_at = ?", id, 0).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func GetCrawlTaskAll() ([]*CrawlTask, error) {
	var tasks []*CrawlTask
	if err := DB.Where("deleted_at = ?", 0).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteCrawlTask(task *CrawlTask) error {
	return DB.Model(task).Where("deleted_at = ?", 0).Update("deleted_at", time.Now().Unix()).Error
}

func ExistCrawlTask(dataType string, dataId int, crawlType string) (bool, error) {
	var task CrawlTask
	err := DB.Select("id").Where("data_type = ? AND data_id = ? AND account_type = ? AND deleted_at = ?",
		dataType, dataId, crawlType, 0).First(&task).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return task.ID > 0, nil
}
