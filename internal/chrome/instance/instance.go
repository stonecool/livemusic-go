package instance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/chrome/util"
	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stonecool/livemusic-go/internal/task"
	"go.uber.org/zap"
)

var _ types.Chrome = (*Instance)(nil)

type Instance struct {
	ID           int
	IP           string
	Port         int
	DebuggerURL  string
	State        types.InstanceState
	mu           sync.RWMutex
	Accounts     map[string]account.IAccount
	StateChan    chan types.StateEvent
	allocatorCtx context.Context
	cancelFunc   context.CancelFunc
	Opts         *types.InstanceOptions
}

func (i *Instance) GetID() int {
	return i.ID
}

func (i *Instance) GetAddr() string {
	return fmt.Sprintf("%s:%d", i.IP, i.Port)
}

func (i *Instance) initialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), i.Opts.InitTimeout)
	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(ctx, i.DebuggerURL)

	i.allocatorCtx = allocatorCtx
	i.cancelFunc = func() {
		allocatorCancel()
		cancel()
	}

	go i.stateManager()
	go i.heartBeat()
	go i.cleanupTabs()

	return nil
}

func (i *Instance) Close() error {
	if i.cancelFunc != nil {
		defer time.Sleep(time.Second)
		i.cancelFunc()

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

func (i *Instance) IsAvailable() bool {
	return i.GetState() == types.InstanceStateAvailable
}

func (i *Instance) isAvailable(cat string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	acc, exists := i.Accounts[cat]

	if !exists {
		return false
	}

	return acc.IsAvailable()
}

func (i *Instance) SetState(state types.InstanceState) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.State = state
}

func (i *Instance) GetState() types.InstanceState {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.State
}

func (i *Instance) GetStateChan() chan types.StateEvent {
	return i.StateChan
}

func (i *Instance) stateManager() {
	for {
		select {
		case evt := <-i.GetStateChan():
			i.HandleStateTransition(evt)
		case <-i.allocatorCtx.Done():
			return
		}
	}
}

func (i *Instance) HandleStateTransition(evt types.StateEvent) {
	oldState := i.GetState()
	var newState types.InstanceState
	var err error

	if !oldState.IsValidTransition(evt.Type) {
		err = fmt.Errorf("invalid state transition from %s with event %v",
			oldState.String(), evt.Type)
		evt.Response <- err
		return
	}

	switch evt.Type {
	case types.EventHealthCheckSuccess:
		newState = types.InstanceStateAvailable
	case types.EventHealthCheckFail:
		switch oldState {
		case types.InstanceStateAvailable:
			newState = types.InstanceStateUnstable
		case types.InstanceStateUnstable:
			// 检查失败次数
			failCount := evt.Data.(int)
			if failCount >= 3 {
				newState = types.InstanceStateUnavailable
			} else {
				newState = types.InstanceStateUnstable
			}
		}
	}

	if oldState != newState {
		i.SetState(newState)
		// 如果变为不可用状态,从实例池中移除
		if newState == types.InstanceStateUnavailable {
			// TODO: 通知实例池移除该实例
		}
	}

	evt.Response <- err
}

func (i *Instance) HandleEvent(event types.EventType) {
	i.GetStateChan() <- types.StateEvent{
		Type: event,
	}
}

func (i *Instance) GetNewContext() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(i.allocatorCtx)
}

func (i *Instance) RetryInitialize(maxAttempts int) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := i.initialize(); err == nil {
			return nil
		} else {
			lastErr = err
			internal.Logger.Warn("initialization attempt failed",
				zap.Int("attempt", attempt),
				zap.Int("maxAttempts", maxAttempts),
				zap.Error(err),
				zap.Int("id", i.ID))
			time.Sleep(time.Second * time.Duration(attempt))
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

// 心跳检查
func (i *Instance) heartBeat() {
	ticker := time.NewTicker(i.Opts.HeartbeatInterval)
	defer ticker.Stop()

	failCount := 0

	for {
		select {
		case <-ticker.C:
			ok, _ := util.RetryCheckChromeHealth(i.GetAddr(), 1, 0)
			if !ok {
				failCount++
				i.HandleEvent(types.EventHealthCheckFail)
				i.GetStateChan() <- types.StateEvent{
					Type: types.EventHealthCheckFail,
					Data: failCount,
				}
			} else {
				failCount = 0
				i.GetStateChan() <- types.StateEvent{
					Type: types.EventHealthCheckSuccess,
				}
			}

		case <-i.allocatorCtx.Done():
			return
		}
	}
}

func (i *Instance) GetAccounts() map[string]account.IAccount {
	i.mu.RLock()
	defer i.mu.RUnlock()

	accounts := make(map[string]account.IAccount, len(i.Accounts))
	for k, v := range i.Accounts {
		accounts[k] = v
	}
	return accounts
}

// NeedsReInitialize 判断是否需要重新初始化
func (i *Instance) NeedsReInitialize() bool {
	state := i.GetState()
	return state == types.InstanceStateUnavailable
}

func (i *Instance) cleanupTabs() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//targets, err := chromedp.Targets(i.allocatorCtx)
			//if err != nil {
			//	internal.Logger.Error("failed to get targets",
			//		zap.Error(err),
			//		zap.Int("chromeID", i.ID))
			//	continue
			//}

			//now := time.Now()
			//for _, t := range targets {
			//	// 跳过主页面和空白页
			//	if t.Type == "page" && t.URL != "about:blank" {
			//		if now.Sub(t.LastActivityTime) > 30*time.Minute {
			//			if err := chromedp.CloseTarget(i.allocatorCtx, t.TargetID); err != nil {
			//				internal.Logger.Error("failed to close target",
			//					zap.Error(err),
			//					zap.String("targetID", string(t.TargetID)),
			//					zap.Int("chromeID", i.ID))
			//			}
			//		}
			//	}
			//}

		case <-i.allocatorCtx.Done():
			return
		}
	}
}

func (i *Instance) ExecuteTask(task task.ITask) error {
	i.mu.RLock()
	defer i.mu.RUnlock()
	acc, exists := i.Accounts[task.GetCategory()]

	if !exists {
		internal.Logger.Error("no account found for category",
			zap.String("category", task.GetCategory()),
			zap.Int("chromeID", i.ID))
		return fmt.Errorf("no account found for category: %s", task.GetCategory())
	}

	if !acc.IsAvailable() {
		internal.Logger.Error("account not available",
			zap.String("category", task.GetCategory()),
			zap.Int("chromeID", i.ID))
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

func (i *Instance) SetAccount(category string, acc account.IAccount) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.Accounts[category] = acc
}

func (i *Instance) checkZombieProcess() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !i.IsAvailable() {
				internal.Logger.Warn("chrome instance appears to be zombie",
					zap.Int("chromeID", i.ID),
					zap.String("addr", i.GetAddr()))

				// 尝试重新初始化
				if err := i.RetryInitialize(3); err != nil {
					internal.Logger.Error("failed to reinitialize zombie chrome",
						zap.Error(err),
						zap.Int("chromeID", i.ID))
				}
			}
		case <-i.allocatorCtx.Done():
			return
		}
	}
}

// 实现 Initialize 接口方法
func (i *Instance) Initialize() error {
	return i.RetryInitialize(3)
}

func (i *Instance) GetModelData() *types.Model {
	return &types.Model{
		BaseModel: database.BaseModel{
			ID: i.ID,
		},
		IP:          i.IP,
		Port:        i.Port,
		DebuggerURL: i.DebuggerURL,
		State:       int(i.State),
	}
}
