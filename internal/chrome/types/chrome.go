package types

import (
	"time"

	"github.com/stonecool/livemusic-go/internal/task"
)

// IChrome 定义 Chrome 实例的基础接口
type IChrome interface {
	GetID() int
	GetAddr() string
	IsAvailable() bool
	ExecuteTask(task task.ITask) error
	Close() error
}

// ChromeState 实例状态
type ChromeState uint8

const (
	ChromeStateConnected ChromeState = iota
	ChromeStateDisconnected
	ChromeStateOffline
)

// InstanceOptions 实例配置选项
type InstanceOptions struct {
	HeartbeatInterval   time.Duration
	InitTimeout         time.Duration
	TabCleanupInterval  time.Duration
	TabInactiveTimeout  time.Duration
	ZombieCheckInterval time.Duration
}
