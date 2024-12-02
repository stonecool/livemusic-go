package instance

import (
	"context"
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/chrome/util"
	"github.com/stonecool/livemusic-go/internal/task"
	"go.uber.org/zap"
)

type Chrome struct {
	ID           int
	IP           string
	Port         int
	Accounts     map[string]account.IAccount
	AccountsMu   sync.RWMutex
	DebuggerURL  string
	State        types.ChromeState
	StateChan    chan types.StateEvent
	allocatorCtx context.Context
	cancelFunc   context.CancelFunc
	Opts         *types.InstanceOptions
}

func (c *Chrome) SetState(state types.ChromeState) {
	c.State = state
}

func (c *Chrome) GetState() types.ChromeState {
	return c.State
}

func (c *Chrome) GetStateChan() chan types.StateEvent {
	return c.StateChan
}

func (c *Chrome) GetID() int {
	return c.ID
}

// 确保 Chrome 实现了 Chrome 接口
var _ types.Chrome = (*Chrome)(nil)

func NewChrome(ip string, port int, url string, state types.ChromeState) *Chrome {
	return &Chrome{
		IP:          ip,
		Port:        port,
		DebuggerURL: url,
		State:       state,
	}
}

func (c *Chrome) stateManager() {
	for {
		select {
		case evt := <-c.GetStateChan():
			switch evt.Type {
			case types.EventGetState:
				evt.Response <- c.GetState()
				continue
			default:
				c.HandleStateTransition(evt)
			}
		case <-c.allocatorCtx.Done():
			return
		}
	}
}

func (c *Chrome) HandleStateTransition(evt types.StateEvent) {
	oldState := c.GetState()
	var newState types.ChromeState
	var err error

	if !oldState.IsValidTransition(evt.Type) {
		err = fmt.Errorf("invalid state transition from %s with event %v",
			oldState.String(), evt.Type)
		evt.Response <- err
		return
	}

	switch evt.Type {
	case types.EventHealthCheckFail:
		newState = types.ChromeStateDisconnected
		//case eventShutdown:
		//	newState = stateOffline
	}

	if oldState != newState {
		c.SetState(newState)

		//if err := chrome.UpdateChrome(c); err != nil {
		//	internal.Logger.Error("failed to update chrome state",
		//		zap.Error(err),
		//		zap.Int("chromeID", c.GetID()),
		//		zap.String("oldState", oldState.String()),
		//		zap.String("newState", newState.String()))
		//}
	}

	evt.Response <- err
}

func (c *Chrome) HandleEvent(event types.EventType) error {
	response := make(chan interface{}, 1)
	c.GetStateChan() <- types.StateEvent{
		Type:     event,
		Response: response,
	}
	result := <-response
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (c *Chrome) initialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Opts.InitTimeout)
	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(ctx, c.DebuggerURL)

	c.allocatorCtx = allocatorCtx
	c.cancelFunc = func() {
		allocatorCancel()
		cancel()
	}

	go c.stateManager()
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
	ticker := time.NewTicker(c.Opts.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			state := c.getState()
			if state != types.ChromeStateConnected {
				continue
			}

			ok, _ := util.RetryCheckChromeHealth(c.GetAddr(), 1, 0)
			if !ok {
				c.HandleEvent(types.EventHealthCheckFail)
			}

		case <-c.allocatorCtx.Done():
			return
		}
	}
}

func (c *Chrome) GetAccounts() map[string]account.IAccount {
	c.AccountsMu.RLock()
	defer c.AccountsMu.RUnlock()

	accounts := make(map[string]account.IAccount, len(c.Accounts))
	for k, v := range c.Accounts {
		accounts[k] = v
	}
	return accounts
}

func (c *Chrome) isAvailable(cat string) bool {
	c.AccountsMu.RLock()
	defer c.AccountsMu.RUnlock()
	acc, exists := c.Accounts[cat]

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
	}
	return nil
}

// 判断实例是否可用
func (c *Chrome) IsAvailable() bool {
	return c.getState() == types.ChromeStateConnected
}

// getState 获取当前状态
func (c *Chrome) getState() types.ChromeState {
	response := make(chan interface{}, 1)
	c.StateChan <- types.StateEvent{
		Type:     types.EventGetState,
		Response: response,
	}
	result := <-response
	return result.(types.ChromeState)
}

// NeedsReInitialize 判断是否需要重新初始化
func (c *Chrome) NeedsReInitialize() bool {
	state := c.getState()
	return state == types.ChromeStateDisconnected
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
	c.AccountsMu.RLock()
	defer c.AccountsMu.RUnlock()
	acc, exists := c.Accounts[task.GetCategory()]

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
	c.AccountsMu.Lock()
	defer c.AccountsMu.Unlock()

	c.Accounts[category] = acc
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

// 实现 Initialize 接口方法
func (c *Chrome) Initialize() error {
	return c.RetryInitialize(3)
}

func (c *Chrome) ToModel() *types.Model {
	//ID           int
	//IP           string
	//Port         int
	//Accounts     map[string]account.IAccount
	//AccountsMu   sync.RWMutex
	//DebuggerURL  string
	//State        types.ChromeState
	//StateChan    chan types.StateEvent
	//allocatorCtx context.Context
	//cancelFunc   context.CancelFunc
	//Opts         *types.InstanceOptions

	baseModel := database.BaseModel{
		ID: c.ID,
	}

	return &types.Model{
		BaseModel:   baseModel,
		IP:          c.IP,
		Port:        c.Port,
		DebuggerURL: c.DebuggerURL,
		State:       int(c.State),
	}
}