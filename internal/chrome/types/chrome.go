package types

import (
	"context"

	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/task"
)

// Chrome 实例的基础接口
type Chrome interface {
	// 基础信息
	GetID() int
	GetAddr() string

	// 生命周期管理
	Initialize() error
	Close() error
	IsAvailable() bool

	// 状态管理
	GetState() ChromeState
	SetState(state ChromeState)
	GetStateChan() chan StateEvent

	// 账号管理
	GetAccounts() map[string]account.IAccount

	// 任务执行
	ExecuteTask(task task.ITask) error
	GetNewContext() (context.Context, context.CancelFunc)
}
