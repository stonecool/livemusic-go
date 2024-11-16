package crawltask

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/stonecool/livemusic-go/internal/config"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateTask(task *CrawlTask) error {
	if err := v.ValidateCategory(task.Category); err != nil {
		return err
	}

	if err := v.ValidateDataType(task.DataType); err != nil {
		return err
	}

	if err := v.ValidateCronSpec(task.CronSpec); err != nil {
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

func (v *Validator) ValidateDataType(dataType string) error {
	if dataType == "" {
		return fmt.Errorf("data type cannot be empty")
	}

	if _, ok := internal.dataType2StructMap[dataType]; !ok {
		return fmt.Errorf("unsupported data type: %s", dataType)
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
