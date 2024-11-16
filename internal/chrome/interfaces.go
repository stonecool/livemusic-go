package chrome

import (
	"context"

	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/task"
)

// IInstance 实例接口
type IInstance interface {
	GetID() int
	GetCategory() string
	GetContext() context.Context
	AddAccount(account account.ICrawlAccount) error
	ExecuteTask(task *task.Task) error
	Close() error
}

// Pool 实例池接口
type Pool interface {
	CreateInstance(category string) (Instance, error)
	GetInstance(id int) (Instance, error)
	GetInstancesByCategory(category string) []Instance
	RemoveInstance(id int) error
	DispatchTask(category string, task *task.Task) error
}

// IFactory 工厂接口
type IFactory interface {
	CreatePool() Pool
}
