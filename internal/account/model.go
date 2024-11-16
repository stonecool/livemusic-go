package account

import (
	"time"

	"github.com/stonecool/livemusic-go/internal"
)

// Model 数据库模型
type Model struct {
	ID          int    `gorm:"primaryKey"`
	Category    string `gorm:"type:varchar(50);not null"` // 账号类型
	AccountName string `gorm:"type:varchar(50);not null"` // 账号名称
	LastURL     string `gorm:"type:varchar(255)"`         // 最后访问的URL
	Cookies     []byte `gorm:"type:bytes"`                // Cookie数据
	InstanceID  int    `gorm:"default:0"`                 // 关联的Chrome实例ID
	State       int    `gorm:"default:0"`                 // 账号状态
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

// TableName 指定表名
func (Model) TableName() string {
	return "crawl_accounts"
}

// ToEntity 转换为领域对象
func (m *Model) ToEntity() *Account {
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

// FromEntity 从领域对象更新模型
func (m *Model) FromEntity(account *Account) {
	m.ID = account.ID
	m.Category = account.Category
	m.AccountName = account.AccountName
	m.LastURL = account.lastURL
	m.Cookies = account.cookies
	m.InstanceID = account.InstanceID
	m.State = int(account.State)
}
