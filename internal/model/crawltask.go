package model

import (
	"time"

	"gorm.io/gorm"
)

type CrawlTask struct {
	Model

	Category  string
	TargetId  string
	MetaType  string
	MetaID    int
	Mark      string
	Count     int
	CronSpec  string // cron表达式,如 "0 0 9 * * *"
	FirstTime int
	LastTime  int
}

func AddCrawlTask(data map[string]interface{}) (*CrawlTask, error) {
	task := CrawlTask{
		Category: data["category"].(string),
		TargetId: data["target_id"].(string),
		MetaType: data["meta_type"].(string),
		MetaID:   data["meta_id"].(int),
		CronSpec: data["cron_spec"].(string),
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

func ExistCrawlTask(dataType string, dataId int, category string) (bool, error) {
	var task CrawlTask
	err := DB.Select("id").Where("meta_type = ? AND meta_id = ? AND category = ? AND deleted_at = ?",
		dataType, dataId, category, 0).First(&task).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return task.ID > 0, nil
}

func EditCrawlTask(id int, data map[string]interface{}) (*CrawlTask, error) {
	var task CrawlTask

	if err := DB.Model(&task).Where("id = ? AND deleted_at = ? ", id, 0).Updates(data).Error; err != nil {
		return nil, err
	}

	return &task, nil
}
