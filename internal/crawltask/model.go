package crawltask

import (
	"time"

	"github.com/stonecool/livemusic-go/internal"
)

type Model struct {
	ID              int    `gorm:"primaryKey"`
	DataType        string `gorm:"type:varchar(50);not null"`
	DataID          int    `gorm:"not null"`
	Category        string `gorm:"type:varchar(50);not null"`
	AccountID       int    `gorm:"not null"`
	CronSpec        string `gorm:"type:varchar(50)"`
	State           int    `gorm:"default:0"`
	LastExecuteTime time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `gorm:"index"`
}

func (Model) TableName() string {
	return "crawl_tasks"
}

func (m *Model) ToEntity() *CrawlTask {
	return &CrawlTask{
		ID:              m.ID,
		DataType:        m.DataType,
		DataID:          m.DataID,
		Category:        m.Category,
		AccountID:       m.AccountID,
		CronSpec:        m.CronSpec,
		State:           internal.TaskState(m.State),
		LastExecuteTime: m.LastExecuteTime,
	}
}

func (m *Model) FromEntity(task *CrawlTask) {
	m.ID = task.ID
	m.DataType = task.DataType
	m.DataID = task.DataID
	m.Category = task.Category
	m.AccountID = task.AccountID
	m.CronSpec = task.CronSpec
	m.State = int(task.State)
	m.LastExecuteTime = task.LastExecuteTime
}
