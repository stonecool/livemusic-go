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

func newChrome(ip string, port int, url string, state chromeState) *Chrome {
	return &Chrome{
		IP:          ip,
		Port:        port,
		DebuggerURL: url,
		State:       state,
	}
}

func (c *Chrome) initialize() error {
	if c.State != chromeStateUninitialized {
		return fmt.Errorf("instance in invalid state: %v", c.State)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.opts.InitTimeout)
	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(ctx, c.DebuggerURL)

	c.allocatorCtx = allocatorCtx
	c.cancelFunc = func() {
		allocatorCancel()
		cancel()
	}

	go stateManager(c)

	if ok, _ := RetryCheckChromeHealth(c.getAddr(), 1, 0); !ok {
		c.cancelFunc()
		handleEvent(c, EventInitFail)
		internal.Logger.Error("instance initialization failed: health check failed",
			zap.String("addr", c.getAddr()),
			zap.Int("id", c.ID))
		return fmt.Errorf("instance initialization failed: health check failed")
	}

	if err := handleEvent(c, EventInitSuccess); err != nil {
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
			state := c.getState()
			if state != chromeStateConnected {
				continue
			}

			ok, _ := RetryCheckChromeHealth(c.getAddr(), 1, 0)
			if !ok {
				handleEvent(c, EventHealthCheckFail)
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

func (c *Chrome) getAddr() string {
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
		if err := handleEvent(c, EventInitFail); err != nil {
			internal.Logger.Error("failed to update state on close",
				zap.Error(err),
				zap.Int("chromeID", c.ID))
		}
	}
	return nil
}

// 判断实例是否可用
func (c *Chrome) IsAvailable() bool {
	return c.getState() == chromeStateConnected
}

// getState 获取当前状态
func (c *Chrome) getState() chromeState {
	response := make(chan interface{}, 1)
	c.stateChan <- stateEvent{
		Type:     EventGetState,
		Response: response,
	}
	result := <-response
	return result.(chromeState)
}

// NeedsReInitialize 判断是否需要重新初始化
func (c *Chrome) NeedsReInitialize() bool {
	state := c.getState()
	return state == chromeStateInitFailed || state == chromeStateDisconnected
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
					zap.String("addr", c.getAddr()))

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
