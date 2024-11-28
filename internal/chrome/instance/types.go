package instance

import (
	"time"
)

// 实例配置选项
type InstanceOptions struct {
	HeartbeatInterval   time.Duration
	InitTimeout         time.Duration
	TabCleanupInterval  time.Duration
	TabInactiveTimeout  time.Duration
	ZombieCheckInterval time.Duration
}

// 默认配置
func DefaultOptions() *InstanceOptions {
	return &InstanceOptions{
		HeartbeatInterval:   30 * time.Second,
		InitTimeout:         30 * time.Second,
		TabCleanupInterval:  5 * time.Minute,
		TabInactiveTimeout:  30 * time.Minute,
		ZombieCheckInterval: 10 * time.Minute,
	}
}
