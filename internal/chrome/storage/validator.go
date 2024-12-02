package storage

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateChrome(model *types.model) error {
	if model.IP == "" {
		return fmt.Errorf("IP cannot be empty")
	}

	if model.Port <= 0 {
		return fmt.Errorf("invalid port number")
	}

	if model.DebuggerURL == "" {
		return fmt.Errorf("debugger URL cannot be empty")
	}

	return nil
}
