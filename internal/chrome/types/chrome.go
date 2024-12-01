package types

import (
	"context"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/task"
)

// Chrome 实例的基础接口
type Chrome interface {
	GetID() int
	GetAddr() string
	IsAvailable() bool
	ExecuteTask(task task.ITask) error
	Close() error
	Initialize() error
	GetStateChan() chan StateEvent
	GetState() ChromeState
	SetState(state ChromeState)
	GetAccounts() map[string]account.IAccount
	GetNewContext() (context.Context, context.CancelFunc)
}
