package chrome

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/task"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account"
	"go.uber.org/zap"
)

type Chrome struct {
	ID           int
	IP           string
	Port         int
	accounts     map[string]account.IAccount
	accountsMu   sync.RWMutex
	DebuggerURL  string
	State        chromeState
	stateChan    chan stateEvent
	allocatorCtx context.Context
	cancelFunc   context.CancelFunc
	opts         *InstanceOptions
}

func (c *Chrome) initialize() error {
	if c.State != STATE_UNINITIALIZED {
		return fmt.Errorf("instance in invalid state: %v", c.State)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.opts.InitTimeout)
	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(ctx, c.DebuggerURL)

	c.allocatorCtx = allocatorCtx
	c.cancelFunc = func() {
		allocatorCancel()
		cancel()
	}

	go c.stateManager()

	if ok, _ := RetryCheckChromeHealth(c.GetAddr(), 1, 0); !ok {
		c.cancelFunc()
		c.handleEvent(EVENT_INIT_FAIL)
		internal.Logger.Error("instance initialization failed: health check failed",
			zap.String("addr", c.GetAddr()),
			zap.Int("id", c.ID))
		return fmt.Errorf("instance initialization failed: health check failed")
	}

	if err := c.handleEvent(EVENT_INIT_SUCCESS); err != nil {
		c.cancelFunc()
		internal.Logger.Error("failed to update state",
			zap.Error(err),
			zap.Int("id", c.ID))
		return fmt.Errorf("failed to update state: %w", err)
	}

	go c.heartBeat()
	go c.cleanupTabs()

	return nil
}

func (c *Chrome) GetNewContext() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(c.allocatorCtx)
}

func (c *Chrome) RetryInitialize(maxAttempts int) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := c.initialize(); err == nil {
			return nil
		} else {
			lastErr = err
			internal.Logger.Warn("initialization attempt failed",
				zap.Int("attempt", attempt),
				zap.Int("maxAttempts", maxAttempts),
				zap.Error(err),
				zap.Int("id", c.ID))
			time.Sleep(time.Second * time.Duration(attempt))
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

// 心跳检查
func (c *Chrome) heartBeat() {
	ticker := time.NewTicker(c.opts.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			state := c.GetState()
			if state != STATE_CONNECTED {
				continue
			}

			ok, _ := RetryCheckChromeHealth(c.GetAddr(), 1, 0)
			if !ok {
				c.handleEvent(EVENT_HEALTH_CHECK_FAIL)
			}

		case <-c.allocatorCtx.Done():
			return
		}
	}
}

func (c *Chrome) getAccounts() map[string]account.IAccount {
	c.accountsMu.RLock()
	defer c.accountsMu.RUnlock()

	accounts := make(map[string]account.IAccount, len(c.accounts))
	for k, v := range c.accounts {
		accounts[k] = v
	}
	return accounts
}

func (c *Chrome) isAvailable(cat string) bool {
	c.accountsMu.RLock()
	defer c.accountsMu.RUnlock()
	acc, exists := c.accounts[cat]

	if !exists {
		return false
	}

	return acc.IsAvailable()
}

func (c *Chrome) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.IP, c.Port)
}

func (c *Chrome) Close() error {
	if c.cancelFunc != nil {
		// 先取消上下文，让所有goroutine优雅退出
		c.cancelFunc()

		// 等待一段时间让goroutine完成清理工作
		time.Sleep(time.Second)

		// 关闭所有打开的标签页
		//targets, err := chromedp.Targets(context.Background())
		//if err == nil {
		//	for _, t := range targets {
		//		if t.Type == "page" {
		//			chromedp.CloseTarget(context.Background(), t.TargetID)
		//		}
		//	}
		//}

		// 更新实例状态
		if err := c.handleEvent(EVENT_INIT_FAIL); err != nil {
			internal.Logger.Error("failed to update state on close",
				zap.Error(err),
				zap.Int("chromeID", c.ID))
		}
	}
	return nil
}

// 判断实例是否可用
func (c *Chrome) IsAvailable() bool {
	return c.GetState() == STATE_CONNECTED
}

// 状态管理器
func (c *Chrome) stateManager() {
	for {
		select {
		case evt := <-c.stateChan:
			var err error
			oldState := c.State

			switch evt.event {
			case EVENT_GET_STATE:
				evt.response <- c.State
				continue

			case EVENT_INIT_SUCCESS:
				if c.State == STATE_UNINITIALIZED {
					c.State = STATE_CONNECTED
				} else {
					err = fmt.Errorf("cannot initialize from state: %v", c.State)
				}

			case EVENT_INIT_FAIL:
				if c.State == STATE_UNINITIALIZED {
					c.State = STATE_INIT_FAILED
				} else {
					err = fmt.Errorf("cannot fail initialization from state: %v", c.State)
				}

			case EVENT_HEALTH_CHECK_SUCCESS:
				if c.State == STATE_DISCONNECTED {
					c.State = STATE_CONNECTED
				}

			case EVENT_HEALTH_CHECK_FAIL:
				if c.State == STATE_CONNECTED {
					c.State = STATE_DISCONNECTED
				}
			}

			if err == nil && oldState != c.State {
				// TODO: 更新实例状态
			}

			evt.response <- err
		case <-c.allocatorCtx.Done():
			return
		}
	}
}

// 处理状态事件
func (c *Chrome) handleEvent(event InstanceEvent) error {
	response := make(chan interface{}, 1)
	c.stateChan <- stateEvent{
		event:    event,
		response: response,
	}
	result := <-response
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

// GetState 获取当前状态
func (c *Chrome) GetState() chromeState {
	response := make(chan interface{}, 1)
	c.stateChan <- stateEvent{
		event:    EVENT_GET_STATE,
		response: response,
	}
	result := <-response
	return result.(chromeState)
}

// NeedsReInitialize 判断是否需要重新初始化
func (c *Chrome) NeedsReInitialize() bool {
	state := c.GetState()
	return state == STATE_INIT_FAILED || state == STATE_DISCONNECTED
}

func (c *Chrome) cleanupTabs() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//targets, err := chromedp.Targets(c.allocatorCtx)
			//if err != nil {
			//	internal.Logger.Error("failed to get targets",
			//		zap.Error(err),
			//		zap.Int("chromeID", c.ID))
			//	continue
			//}

			//now := time.Now()
			//for _, t := range targets {
			//	// 跳过主页面和空白页
			//	if t.Type == "page" && t.URL != "about:blank" {
			//		if now.Sub(t.LastActivityTime) > 30*time.Minute {
			//			if err := chromedp.CloseTarget(c.allocatorCtx, t.TargetID); err != nil {
			//				internal.Logger.Error("failed to close target",
			//					zap.Error(err),
			//					zap.String("targetID", string(t.TargetID)),
			//					zap.Int("chromeID", c.ID))
			//			}
			//		}
			//	}
			//}

		case <-c.allocatorCtx.Done():
			return
		}
	}
}

func (c *Chrome) ExecuteTask(task task.ITask) error {
	c.accountsMu.RLock()
	defer c.accountsMu.RUnlock()
	acc, exists := c.accounts[task.GetCategory()]

	if !exists {
		internal.Logger.Error("no account found for category",
			zap.String("category", task.GetCategory()),
			zap.Int("chromeID", c.ID))
		return fmt.Errorf("no account found for category: %s", task.GetCategory())
	}

	if !acc.IsAvailable() {
		internal.Logger.Error("account not available",
			zap.String("category", task.GetCategory()),
			zap.Int("chromeID", c.ID))
		return fmt.Errorf("account not available")
	}

	// TODO
	//select {
	//case acc.TaskChan <- task:
	//	return nil
	//case <-time.After(5 * time.Second):
	//	return fmt.Errorf("send task timeout")
	//}
	return nil
}

func (c *Chrome) SetAccount(category string, acc account.IAccount) {
	c.accountsMu.Lock()
	defer c.accountsMu.Unlock()

	c.accounts[category] = acc
}

func (c *Chrome) checkZombieProcess() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !c.IsAvailable() {
				internal.Logger.Warn("chrome instance appears to be zombie",
					zap.Int("chromeID", c.ID),
					zap.String("addr", c.GetAddr()))

				// 尝试重新初始化
				if err := c.RetryInitialize(3); err != nil {
					internal.Logger.Error("failed to reinitialize zombie chrome",
						zap.Error(err),
						zap.Int("chromeID", c.ID))
				}
			}
		case <-c.allocatorCtx.Done():
			return
		}
	}
}
