package task

// ITask 基础任务接口
type ITask interface {
	GetID() int
	GetCategory() string
	GetCronSpec() string
}
