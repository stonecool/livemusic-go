package task

import (
	"github.com/stonecool/livemusic-go/internal/database"
)

type model struct {
	database.BaseModel

	Category  string `gorm:"type:varchar(50);not null"`
	TargetID  string `gorm:"not null"`
	MetaType  string `gorm:"type:varchar(50);not null"`
	MetaID    int    `gorm:"not null"`
	Count     int
	CronSpec  string `gorm:"type:varchar(50)"`
	FirstTime int    `gorm:"not null"`
	LastTime  int    `gorm:"not null"`
	Mark      string `gorm:"type:varchar(50)"`
}

func (*model) TableName() string {
	return "tasks"
}

func (m *model) ToEntity() *Task {
	return &Task{
		ID:        m.ID,
		Category:  m.Category,
		TargetID:  m.TargetID,
		MetaType:  m.MetaType,
		MetaID:    m.MetaID,
		Count:     m.Count,
		CronSpec:  m.CronSpec,
		FirstTime: m.FirstTime,
		LastTime:  m.LastTime,
		mark:      m.Mark,
	}
}

func (m *model) FromEntity(task *Task) {
	m.ID = task.ID
	m.Category = task.Category
	m.TargetID = task.TargetID
	m.MetaType = task.MetaType
	m.MetaID = task.MetaID
	m.Count = task.Count
	m.CronSpec = task.CronSpec
	m.FirstTime = task.FirstTime
	m.LastTime = task.LastTime
	m.Mark = task.mark
}
