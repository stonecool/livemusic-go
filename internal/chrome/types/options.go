package types

import (
	"fmt"
	"time"
)

// InstanceOptions 实例配置选项
type InstanceOptions struct {
	HeartbeatInterval   time.Duration
	InitTimeout         time.Duration
	TabCleanupInterval  time.Duration
	TabInactiveTimeout  time.Duration
	ZombieCheckInterval time.Duration
}

// Validate 验证配置选项
func (o *InstanceOptions) Validate() error {
	if o.HeartbeatInterval <= 0 {
		return fmt.Errorf("heartbeat interval must be positive")
	}
	if o.InitTimeout <= 0 {
		return fmt.Errorf("init timeout must be positive")
	}
	// ... 其他验证
	return nil
}
