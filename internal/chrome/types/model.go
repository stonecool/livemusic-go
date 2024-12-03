package types

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type Model struct {
	database.BaseModel

	IP          string `gorm:"type:varchar(20);not null"`
	Port        int    `gorm:"default:0"`
	DebuggerURL string `gorm:"type:varchar(100);not null"`
	State       int    `gorm:"default:0"`
}

func (*Model) TableName() string {
	return "chromes"
}

func (m *Model) Validate() error {
	if m.IP == "" {
		return fmt.Errorf("IP cannot be empty")
	}

	if m.Port <= 0 {
		return fmt.Errorf("invalid port number")
	}

	if m.DebuggerURL == "" {
		return fmt.Errorf("debugger URL cannot be empty")
	}

	return nil
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	return m.Validate()
}

func (m *Model) BeforeUpdate(tx *gorm.DB) error {
	return m.Validate()
}

func (m *Model) GetID() int {
	return m.ID
}
