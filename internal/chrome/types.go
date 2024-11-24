package chrome

import (
	"time"
)

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
