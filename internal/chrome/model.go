package chrome

import (
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type model struct {
	database.BaseModel

	IP          string
	Port        int
	DebuggerURL string
	State       int
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
		State:       ChromeState(m.State),
		accounts:    make(map[string]*account.Account),
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
	v := NewValidator()
	return v.ValidateChrome(&Chrome{
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
