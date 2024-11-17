package task

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/config"
)

type Factory struct {
	repo Repository
}

func newFactory(repo Repository) *Factory {
	return &Factory{repo: repo}
}

func (f *Factory) createTask(category string, metaType string, metaId int, cronSpec string) (*Task, error) {
	_, ok := config.AccountMap[category]
	if !ok {
		return nil, fmt.Errorf("account_type:%s not exists",category)
	}

	//exist, err := dataTypeIdExists(metaType, metaId)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if !exist {
	//	return nil, fmt.Errorf("data table not exists")
	//}

	if exist, err := f.repo.ExistsByMeta(metaType, metaId, category); err != nil {
		internal.Logger.Warn("m exists")
		return nil, fmt.Errorf("some error")
	} else if exist {
		return nil, fmt.Errorf("exists")
	}

	v := NewValidator()
	if err := v.ValidateCategory(category); err != nil {
		return nil, fmt.Errorf("invalid task category: %w", err)
	}

	model := &model{
		Category:  category,
		MetaType: metaType,
		MetaID: metaId,
		CronSpec:  cronSpec,
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
