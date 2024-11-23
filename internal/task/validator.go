package task

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"

	"github.com/robfig/cron/v3"
	"github.com/stonecool/livemusic-go/internal/config"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateTask(model *model) error {
	exist, err := repo.existsByMeta(model.Category, model.MetaType, model.MetaID)
	if !exist || err != nil {
		return fmt.Errorf("exists")
	}

	if err := v.ValidateCategory(model.Category); err != nil {
		return err
	}

	if err := v.ValidateMetaType(model.MetaType); err != nil {
		return err
	}

	if err := v.ValidateCronSpec(model.CronSpec); err != nil {
		return err
	}

	return nil
}

func (v *Validator) ValidateCategory(category string) error {
	if category == "" {
		return fmt.Errorf("category cannot be empty")
	}

	if _, ok := config.AccountMap[category]; !ok {
		return fmt.Errorf("unsupported task category: %s", category)
	}

	return nil
}

func (v *Validator) ValidateMetaType(metaType string) error {
	if metaType == "" {
		return fmt.Errorf("data type cannot be empty")
	}

	if _, ok := internal.DataType2StructMap[metaType]; !ok {
		return fmt.Errorf("unsupported data type: %s", metaType)
	}

	return nil
}

func (v *Validator) ValidateCronSpec(spec string) error {
	if spec == "" {
		return fmt.Errorf("cron spec cannot be empty")
	}

	// 验证cron表达式格式
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow)
	if _, err := parser.Parse(spec); err != nil {
		return fmt.Errorf("invalid cron spec: %w", err)
	}

	return nil
}
