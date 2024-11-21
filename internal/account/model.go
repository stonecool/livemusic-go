package account

import (
	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stonecool/livemusic-go/internal/message"
	"gorm.io/gorm"
)

type accountModel struct {
	database.BaseModel

	Category    string `gorm:"type:varchar(50);not null"`
	AccountName string `gorm:"type:varchar(50);not null"`
	LastURL     string `gorm:"type:varchar(255)"`
	Cookies     []byte `gorm:"type:bytes"`
	InstanceID  int    `gorm:"default:0"`
	State       int    `gorm:"default:0"`
}

func (*accountModel) TableName() string {
	return "accounts"
}

func (m *accountModel) toEntity() *Account {
	return &Account{
		ID:          m.ID,
		Category:    m.Category,
		AccountName: m.AccountName,
		lastURL:     m.LastURL,
		cookies:     m.Cookies,
		InstanceID:  m.InstanceID,
		State:       state(m.State),
		msgChan:     make(chan *message.AsyncMessage),
		done:        make(chan struct{}),
	}
}

func (m *accountModel) fromEntity(account *Account) {
	m.ID = account.ID
	m.Category = account.Category
	m.AccountName = account.AccountName
	m.LastURL = account.lastURL
	m.Cookies = account.cookies
	m.InstanceID = account.InstanceID
	m.State = int(account.State)
}

func (m *accountModel) Validate() error {
	v := NewValidator()
	return v.ValidateAccount(&Account{
		Category:    m.Category,
		AccountName: m.AccountName,
	})
}

func (m *accountModel) BeforeCreate(tx *gorm.DB) error {
	return m.Validate()
}

func (m *accountModel) BeforeUpdate(tx *gorm.DB) error {
	return m.Validate()
}

func (m *accountModel) GetID() int {
	return m.ID
}
