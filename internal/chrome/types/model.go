package types

import (
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/chrome/storage"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
	"sync"
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

func (m *Model) ToEntity() *instance.Chrome {
	return &instance.Chrome{
		ID:          m.ID,
		IP:          m.IP,
		Port:        m.Port,
		DebuggerURL: m.DebuggerURL,
		State:       ChromeState(m.State),
		Accounts:    make(map[string]account.IAccount),
		AccountsMu:  sync.RWMutex{},
		StateChan:   make(chan StateEvent),
		Opts:        instance.DefaultOptions(),
	}
}

func (m *Model) FromEntity(chrome *instance.Chrome) {
	m.ID = chrome.ID
	m.IP = chrome.IP
	m.Port = chrome.Port
	m.DebuggerURL = chrome.DebuggerURL
	m.State = int(chrome.State)
}

func (m *Model) Validate() error {
	return storage.NewValidator().ValidateChrome(m)
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
