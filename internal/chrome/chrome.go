package chrome

import (
	"context"
	"fmt"
	"time"

	"github.com/stonecool/livemusic-go/internal/task"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account"
)

type Chrome struct {
	ID           int
	IP           string
	Port         int
	accounts     map[string]account.IAccount
	DebuggerURL  string
	State        State
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
		return fmt.Errorf("instance initialization failed: health check failed")
	}

	if err := c.handleEvent(EVENT_INIT_SUCCESS); err != nil {
		c.cancelFunc()
		return fmt.Errorf("failed to update state: %w", err)
	}

	go c.heartBeat()

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
			time.Sleep(time.Second * time.Duration(attempt)) // 指数退避
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
	return c.accounts
}

func (c *Chrome) isAvailable(cat string) bool {
	account, exists := c.accounts[cat]
	if !exists {
		return false
	}

	return account.IsAvailable()
}

func (c *Chrome) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.IP, c.Port)
}

func (c *Chrome) Close() error {
	if c.cancelFunc != nil {
		c.cancelFunc() // 取消 context，会导致 stateManager 和 heartBeat 退出
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
func (c *Chrome) GetState() State {
	response := make(chan interface{}, 1)
	c.stateChan <- stateEvent{
		event:    EVENT_GET_STATE,
		response: response,
	}
	result := <-response
	return result.(State)
}

// NeedsReInitialize 判断是否需要重新初始化
func (c *Chrome) NeedsReInitialize() bool {
	state := c.GetState()
	return state == STATE_INIT_FAILED || state == STATE_DISCONNECTED
}

func (c *Chrome) cleanupTabs() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
	defer ticker.Stop()

	for {
		select {
		//case <-ticker.C:
		//	// 获取所有 targets (tabs)
		//	targets, err := chromedp.Targets(c.allocatorCtx)
		//	if err != nil {
		//		continue
		//	}
		//
		//	now := time.Now()
		//	for _, t := range targets {
		//		// 跳过主页面
		//		if t.Type == "page" && t.URL != "about:blank" {
		//			// 如果 tab 超过30分钟没有活动，关闭它
		//			if now.Sub(t.LastActivityTime) > 30*time.Minute {
		//				chromedp.CloseTarget(c.allocatorCtx, t.TargetID)
		//			}
		//		}
		//	}

		case <-c.allocatorCtx.Done():
			return
		}
	}
}

func (c *Chrome) ExecuteTask(task task.ITask) error {
	account, exists := c.accounts[task.GetCategory()]
	if !exists {
		return fmt.Errorf("no account found for category: %s", task.GetCategory())
	}

	if !account.IsAvailable() {
		return fmt.Errorf("account not available")
	}

	//select {
	//case account.TaskChan <- task:
	//	return nil
	//case <-time.After(5 * time.Second):
	//	return fmt.Errorf("send task timeout")
	//}
	return nil
}
