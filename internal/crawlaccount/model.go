package crawlaccount

import (
	"github.com/stonecool/livemusic-go/internal"
)

type accountModel struct {
	internal.RawModel

	Category    string `gorm:"type:varchar(50);not null"`
	AccountName string `gorm:"type:varchar(50);not null"`
	LastURL     string `gorm:"type:varchar(255)"`
	Cookies     []byte `gorm:"type:bytes"`
	InstanceID  int    `gorm:"default:0"`
	State       int    `gorm:"default:0"`
}

func (*accountModel) TableName() string {
	return "crawl_accounts"
}

func (m *accountModel) toEntity() *Account {
	return &Account{
		ID:          m.ID,
		Category:    m.Category,
		AccountName: m.AccountName,
		lastURL:     m.LastURL,
		cookies:     m.Cookies,
		InstanceID:  m.InstanceID,
		State:       internal.AccountState(m.State),
		msgChan:     make(chan *internal.AsyncMessage),
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
