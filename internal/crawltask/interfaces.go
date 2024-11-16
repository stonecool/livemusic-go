package crawltask

import (
	"github.com/stonecool/livemusic-go/internal"
)

// Task 基础任务接口
type Task interface {
	GetID() int
	GetCategory() string
	GetState() internal.TaskState
	SetState(state internal.TaskState)
	Execute() error
	Cancel() error
}

// CronTask 定时任务接口
type CronTask interface {
	Task
	GetCronSpec() string
	SetCronSpec(spec string) error
}

// Repository 仓储接口
type Repository interface {
	Create(task *CrawlTask) error
	Get(id int) (*CrawlTask, error)
	Update(task *CrawlTask) error
	Delete(id int) error
	GetAll() ([]*CrawlTask, error)
	FindByCategory(category string) ([]*CrawlTask, error)
}

// Factory 工厂接口
type Factory interface {
	CreateTask(category string, data map[string]interface{}) (*CrawlTask, error)
}
