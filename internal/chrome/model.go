package chrome

import (
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type model struct {
	database.BaseModel

	IP          string `gorm:"type:varchar(20);not null"`
	Port        int    `gorm:"default:0"`
	DebuggerURL string `gorm:"type:varchar(100);not null"`
	State       int    `gorm:"default:0"`
}

func (*model) TableName() string {
	return "chromes"
}

func (m *model) toEntity() *Chrome {
	return &Chrome{
		ID:          m.ID,
		IP:          m.IP,
		Port:        m.Port,
		DebuggerURL: m.DebuggerURL,
		State:       State(m.State),
		accounts:    make(map[string]account.IAccount),
		stateChan:   make(chan stateEvent),
		opts:        DefaultOptions(),
	}
}

func (m *model) fromEntity(chrome *Chrome) {
	m.ID = chrome.ID
	m.IP = chrome.IP
	m.Port = chrome.Port
	m.DebuggerURL = chrome.DebuggerURL
	m.State = int(chrome.State)
}

func (m *model) Validate() error {
	return NewValidator().ValidateChrome(&Chrome{
		IP:   m.IP,
		Port: m.Port,
	})
}

func (m *model) BeforeCreate(tx *gorm.DB) error {
	return m.Validate()
}

func (m *model) BeforeUpdate(tx *gorm.DB) error {
	return m.Validate()
}

func (m *model) GetID() int {
	return m.ID
}
