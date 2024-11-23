package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/config"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateAccount(account *account) error {
	if err := v.validateCategory(account.Category); err != nil {
		return err
	}

	if err := v.validateAccountName(account.GetName()); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateCategory(category string) error {
	if category == "" {
		return fmt.Errorf("category cannot be empty")
	}

	// 检查是否是支持的账号类型
	if _, ok := config.AccountMap[category]; !ok {
		return fmt.Errorf("unsupported account category: %s", category)
	}

	return nil
}

func (v *Validator) validateAccountName(name string) error {
	// 账号名称长度限制
	if len(name) > 50 {
		return fmt.Errorf("account name too long (max 50 characters)")
	}

	return nil
}
