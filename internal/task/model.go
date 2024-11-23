package task

import (
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type model struct {
	database.BaseModel

	Category  string `gorm:"type:varchar(50);not null"`
	TargetId  string `gorm:"type:varchar(100);not null"`
	MetaType  string `gorm:"type:varchar(50);not null"`
	MetaId    int    `gorm:"not null"`
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
		TargetId:  m.TargetId,
		MetaType:  m.MetaType,
		MetaId:    m.MetaId,
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
	m.TargetId = task.TargetId
	m.MetaType = task.MetaType
	m.MetaId = task.MetaId
	m.Mark = task.mark
	m.Count = task.Count
	m.CronSpec = task.CronSpec
	m.FirstTime = task.FirstTime
	m.LastTime = task.LastTime
}

func (m *model) Validate() error {
	//v := NewValidator()
	//return v.ValidateAccount(&account{
	//	Category: m.Category,
	//	Name:     m.Name,
	//})

	return nil
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
