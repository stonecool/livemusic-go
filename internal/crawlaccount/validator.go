package crawlaccount

import (
	"fmt"
	"regexp"

	"github.com/stonecool/livemusic-go/internal/config"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateAccount(account *Account) error {
	if err := v.validateCategory(account.GetCategory()); err != nil {
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
		return fmt.Errorf("unsupported crawlaccount category: %s", category)
	}

	return nil
}

func (v *Validator) validateAccountName(name string) error {
	if name == "" {
		return fmt.Errorf("crawlaccount name cannot be empty")
	}

	// 账号名称长度限制
	if len(name) > 50 {
		return fmt.Errorf("crawlaccount name too long (max 50 characters)")
	}

	// 账号名称格式验证(只允许字母数字下划线)
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, name)
	if !matched {
		return fmt.Errorf("invalid crawlaccount name format")
	}

	return nil
}
