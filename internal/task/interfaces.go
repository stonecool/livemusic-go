package task

// ITask 基础任务接口
type ITask interface {
	GetID() int
	GetCategory() string
	Execute() error
	Cancel() error
}

// CronTask 定时任务接口
type CronTask interface {
	GetCronSpec() string
	SetCronSpec(spec string) error
}

type Repository interface {
	Create(task *Task) error
	Get(id int) (*Task, error)
	Update(task *Task) error
	Delete(id int) error
	GetAll() ([]*Task, error)
	FindByCategory(category string) ([]*Task, error)
	ExistsByMeta(metaType string, metaID int, category string) (bool, error)
}
