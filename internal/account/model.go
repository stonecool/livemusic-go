package account

import (
	"github.com/stonecool/livemusic-go/internal"
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

func (am *accountModel) toEntity() *Account {
	return &Account{
		ID:          am.ID,
		Category:    am.Category,
		AccountName: am.AccountName,
		lastURL:     am.LastURL,
		cookies:     am.Cookies,
		InstanceID:  am.InstanceID,
		State:       internal.AccountState(am.State),
		msgChan:     make(chan *message.AsyncMessage),
		done:        make(chan struct{}),
	}
}

func (am *accountModel) fromEntity(account *Account) {
	am.ID = account.ID
	am.Category = account.Category
	am.AccountName = account.AccountName
	am.LastURL = account.lastURL
	am.Cookies = account.cookies
	am.InstanceID = account.InstanceID
	am.State = int(account.State)
}

func (am *accountModel) Validate() error {
	v := NewValidator()
	return v.ValidateAccount(&Account{
		Category:    am.Category,
		AccountName: am.AccountName,
	})
}

func (am *accountModel) BeforeCreate(tx *gorm.DB) error {
	return am.Validate()
}

func (am *accountModel) BeforeUpdate(tx *gorm.DB) error {
	return am.Validate()
}

func (am *accountModel) GetID() int {
	return am.ID
}