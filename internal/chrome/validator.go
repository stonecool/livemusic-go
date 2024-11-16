package chrome

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/config"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateInstance(instance *Instance) error {
	if err := v.ValidateCategory(instance.Category); err != nil {
		return err
	}

	if instance.AllocatorCtx == nil {
		return fmt.Errorf("allocator context cannot be nil")
	}

	return nil
}

func (v *Validator) ValidateCategory(category string) error {
	if category == "" {
		return fmt.Errorf("category cannot be empty")
	}

	// 检查是否是支持的账号类型
	if _, ok := config.AccountMap[category]; !ok {
		return fmt.Errorf("unsupported instance category: %s", category)
	}

	return nil
}
