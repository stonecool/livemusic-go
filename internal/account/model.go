package account

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/client"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type model struct {
	database.BaseModel

	Category    string `gorm:"type:varchar(50);not null"`
	AccountName string `gorm:"type:varchar(50);not null"`
	LastURL     string `gorm:"type:varchar(255)"`
	Cookies     []byte `gorm:"type:bytes"`
	InstanceID  int    `gorm:"default:0"`
	State       int    `gorm:"default:0"`
}

func (*model) TableName() string {
	return "accounts"
}

func (m *model) toEntity() *Account {
	return &Account{
		ID:          m.ID,
		Category:    m.Category,
		AccountName: m.AccountName,
		lastURL:     m.LastURL,
		cookies:     m.Cookies,
		InstanceID:  m.InstanceID,
		State:       internal.AccountState(m.State),
		msgChan:     make(chan *client.AsyncMessage),
		done:        make(chan struct{}),
	}
}

func (m *model) fromEntity(account *Account) {
	m.ID = account.ID
	m.Category = account.Category
	m.AccountName = account.AccountName
	m.LastURL = account.lastURL
	m.Cookies = account.cookies
	m.InstanceID = account.InstanceID
	m.State = int(account.State)
}

func (m *model) Validate() error {
	v := NewValidator()
	return v.ValidateAccount(&Account{
		Category:    m.Category,
		AccountName: m.AccountName,
	})
}

func (m *model) BeforeCreate(tx *gorm.DB) error {
	return m.Validate()
}

func (m *model) BeforeUpdate(tx *gorm.DB) error {
	return m.Validate()
}