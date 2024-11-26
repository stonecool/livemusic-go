package account

import (
	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stonecool/livemusic-go/internal/message"
	"gorm.io/gorm"
)

type model struct {
	database.BaseModel

	Category   string `gorm:"type:varchar(50);not null"`
	Name       string `gorm:"type:varchar(255);not null"`
	LastURL    string `gorm:"type:varchar(255)"`
	Cookies    []byte `gorm:"type:bytes"`
	InstanceID int    `gorm:"default:0"`
	State      int    `gorm:"default:0"`
}

func (*model) TableName() string {
	return "accounts"
}

func (m *model) toEntity() *account {
	return &account{
		ID:         m.ID,
		Category:   m.Category,
		Name:       m.Name,
		lastURL:    m.LastURL,
		cookies:    m.Cookies,
		instanceID: m.InstanceID,
		State:      accountState(m.State),
		msgChan:    make(chan *message.AsyncMessage),
		done:       make(chan struct{}),
	}
}

func (m *model) fromEntity(account *account) {
	m.ID = account.ID
	m.Category = account.Category
	m.Name = account.Name
	m.LastURL = account.lastURL
	m.Cookies = account.cookies
	m.InstanceID = account.instanceID
	m.State = int(account.State)
}

func (m *model) Validate() error {
	return NewValidator().validateModel(m)
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
