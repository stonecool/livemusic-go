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
	Initialize()
	IsAvailable() bool
	Close()

	// 状态管理
	GetState() InstanceState
	SetState(InstanceState)

	// 账号管理
	GetAccounts() map[string]account.IAccount

	// 任务执行
	ExecuteTask(task.ITask) error
	GetNewContext() (context.Context, context.CancelFunc)
}

type InstanceType uint8

const (
	InstanceTypePersistent InstanceType = iota
	InstanceTypeTemporary
)

func (t InstanceType) String() string {
	switch t {
	case InstanceTypePersistent:
		return "Persistent"
	case InstanceTypeTemporary:
		return "Temporary"
	default:
		return "Unknown"
	}
}
