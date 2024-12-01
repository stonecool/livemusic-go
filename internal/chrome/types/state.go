package types

import (
	"time"
)

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

// String 返回状态的字符串表示
func (s ChromeState) String() string {
	switch s {
	case ChromeStateConnected:
		return "Connected"
	case ChromeStateDisconnected:
		return "Disconnected"
	case ChromeStateOffline:
		return "Offline"
	default:
		return "Unknown"
	}
}

func (s ChromeState) IsValidTransition(event EventType) bool {
	switch event {
	case EventHealthCheckFail:
		return s == ChromeStateConnected
	default:
		return true
	}
}

type EventType uint8

const (
	EventHealthCheckFail EventType = iota
	EventGetState
)
