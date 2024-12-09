package types

import (
	"context"

	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/task"
)

type InstanceType uint8

const (
	InstanceTypePersistent InstanceType = iota
	InstanceTypeTemporary
)

func (t InstanceType) String() string {
	switch t {
	case InstanceTypePersistent:
		return "Persistent" // 持久化实例
	case InstanceTypeTemporary:
		return "Temporary" // 临时实例
	default:
		return "Unknown"
	}
}

// Chrome 实例的基础接口
type Chrome interface {
	// 基础信息
	GetID() int
	GetAddr() string
	GetType() InstanceType // 新增：获取实例类型
	SetType(InstanceType)

	// 生命周期管理
	Initialize() error
	Close() error
	IsAvailable() bool

	// 状态管理
	GetState() InstanceState
	SetState(InstanceState)
	GetStateChan() chan StateEvent

	// 账号管理
	GetAccounts() map[string]account.IAccount

	// 任务执行
	ExecuteTask(task.ITask) error
	GetNewContext() (context.Context, context.CancelFunc)
}
