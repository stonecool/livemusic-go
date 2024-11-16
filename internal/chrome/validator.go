package chrome

import (
	"fmt"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateChrome(chrome *Chrome) error {
	// 验证基本属性
	if chrome.IP == "" {
		return fmt.Errorf("IP cannot be empty")
	}

	if chrome.Port <= 0 {
		return fmt.Errorf("invalid port number")
	}

	if chrome.DebuggerURL == "" {
		return fmt.Errorf("debugger URL cannot be empty")
	}

	return nil
}
