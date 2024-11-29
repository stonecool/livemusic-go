package types

import (
	"github.com/stonecool/livemusic-go/internal/task"
	"time"
)

// Chrome 实例的基础接口
type IChrome interface {
	GetID() int
	GetAddr() string
	IsAvailable() bool
	ExecuteTask(task task.ITask) error
	Close() error
}

// Chrome 实例状态
type ChromeState uint8

const (
	ChromeStateConnected ChromeState = iota
	ChromeStateDisconnected
	ChromeStateOffline
)

// 实例配置选项
type InstanceOptions struct {
	HeartbeatInterval   time.Duration
	InitTimeout         time.Duration
	TabCleanupInterval  time.Duration
	TabInactiveTimeout  time.Duration
	ZombieCheckInterval time.Duration
}

// StateEvent 状态事件
type StateEvent struct {
	Type     EventType
	Response chan interface{}
}

type EventType uint8

const (
	EventHealthCheckFail EventType = iota
	EventGetState
)
