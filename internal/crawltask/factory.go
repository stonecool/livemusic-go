package crawltask

import (
	"fmt"
)

type factoryImpl struct {
	repo Repository
}

func NewFactory(repo Repository) Factory {
	return &factoryImpl{repo: repo}
}

func (f *factoryImpl) CreateTask(category string, data map[string]interface{}) (*CrawlTask, error) {
	v := NewValidator()
	if err := v.ValidateCategory(category); err != nil {
		return nil, fmt.Errorf("invalid task category: %w", err)
	}

	model := &Model{
		Category:  category,
		DataType:  data["data_type"].(string),
		DataID:    data["data_id"].(int),
		AccountID: data["account_id"].(int),
		CronSpec:  data["cron_spec"].(string),
	}

	task := model.ToEntity()
	if err := v.ValidateTask(task); err != nil {
		return nil, fmt.Errorf("invalid task: %w", err)
	}

	if err := f.repo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}
