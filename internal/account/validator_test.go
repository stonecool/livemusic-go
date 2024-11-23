package account

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidator_ValidateCategory(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{"empty category", "", true},
		{"invalid category", "invalid", true},
		{"valid category", "wechat", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateCategory(tt.category)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateAccountName(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name    string
		accName string
		wantErr bool
	}{
		{"empty name", "", false},
		{"too long name", strings.Repeat("a", 256), true},
		{"valid name", "test_123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateAccountName(tt.accName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
