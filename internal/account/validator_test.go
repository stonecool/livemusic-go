package account

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		{"empty name", "", true},
		{"too long name", strings.Repeat("a", 51), true},
		{"invalid chars", "test@123", true},
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
