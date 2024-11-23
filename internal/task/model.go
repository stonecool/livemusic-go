package task

import (
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type model struct {
	database.BaseModel

	Category  string `gorm:"type:varchar(50);not null"`
	TargetID  string `gorm:"type:varchar(100);not null"`
	MetaType  string `gorm:"type:varchar(50);not null"`
	MetaID    int    `gorm:"not null"`
	Mark      string `gorm:"type:varchar(50)"`
	CronSpec  string `gorm:"type:varchar(50)"`
	FirstTime int    `gorm:"not null"`
	LastTime  int    `gorm:"not null"`
	Count     int
}

func (*model) TableName() string {
	return "tasks"
}

func (m *model) toEntity() *Task {
	return &Task{
		ID:        m.ID,
		Category:  m.Category,
		TargetID:  m.TargetID,
		MetaType:  m.MetaType,
		MetaID:    m.MetaID,
		mark:      m.Mark,
		Count:     m.Count,
		CronSpec:  m.CronSpec,
		FirstTime: m.FirstTime,
		LastTime:  m.LastTime,
	}
}

func (m *model) fromEntity(task *Task) {
	m.ID = task.ID
	m.Category = task.Category
	m.TargetID = task.TargetID
	m.MetaType = task.MetaType
	m.MetaID = task.MetaID
	m.Mark = task.mark
	m.Count = task.Count
	m.CronSpec = task.CronSpec
	m.FirstTime = task.FirstTime
	m.LastTime = task.LastTime
}

func (m *model) Validate() error {
	return NewValidator().ValidateTask(m)
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
