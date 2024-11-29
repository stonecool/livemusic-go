package storage

import (
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
	"sync"
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

func (m *model) toEntity() *instance.Chrome {
	return &instance.Chrome{
		ID:          m.ID,
		IP:          m.IP,
		Port:        m.Port,
		DebuggerURL: m.DebuggerURL,
		State:       types.ChromeState(m.State),
		Accounts:    make(map[string]account.IAccount),
		AccountsMu:  sync.RWMutex{},
		StateChan:   make(chan types.StateEvent),
		Opts:        instance.DefaultOptions(),
	}
}

func (m *model) fromEntity(chrome *instance.Chrome) {
	m.ID = chrome.ID
	m.IP = chrome.IP
	m.Port = chrome.Port
	m.DebuggerURL = chrome.DebuggerURL
	m.State = int(chrome.State)
}

func (m *model) Validate() error {
	return NewValidator().ValidateChrome(m)
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
