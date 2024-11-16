package chrome

import (
	"time"
)

// 实例状态
type ChromeState uint8

const (
	STATE_UNINITIALIZED ChromeState = iota // 未初始化：实例刚创建
	STATE_INIT_FAILED                      // 初始化失败：初始化方法执行失败
	STATE_CONNECTED                        // 连接成功：包含初始化成功和心跳检查正常
	STATE_DISCONNECTED                     // 连接断开：心跳检查失败
)

// 实例事件
type InstanceEvent uint8

const (
	EVENT_START                InstanceEvent = iota // 开始初始化
	EVENT_INIT_SUCCESS                              // 初始化成功
	EVENT_INIT_FAIL                                 // 初始化失败
	EVENT_HEALTH_CHECK_SUCCESS                      // 心跳检查成功
	EVENT_HEALTH_CHECK_FAIL                         // 心跳检查失败
	EVENT_GET_STATE                                 // 获取状态
)

// 状态事件
type stateEvent struct {
	event    InstanceEvent
	response chan interface{}
}

// 实例配置选项
type InstanceOptions struct {
	HeartbeatInterval time.Duration
	InitTimeout       time.Duration
}

// 默认配置
func DefaultOptions() *InstanceOptions {
	return &InstanceOptions{
		HeartbeatInterval: 30 * time.Second,
		InitTimeout:       30 * time.Second,
	}
}
