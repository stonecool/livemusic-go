package chrome

import (
	"context"
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"time"

	"github.com/stonecool/livemusic-go/internal/task"

	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/cache"
)

type Chrome struct {
	ID           int
	IP           string
	Port         int
	accounts     map[string]*account.Account
	DebuggerURL  string
	State        ChromeState
	stateChan    chan stateEvent
	allocatorCtx context.Context
	cancelFunc   context.CancelFunc
	opts         *InstanceOptions
}

var instanceCache *cache.Memo

func init() {
	instanceCache = cache.New(getInstance)
}

func getInstance(id int) (interface{}, error) {
	repo := NewRepositoryDB(database.DB)
	factory := NewFactory(repo)
	return factory.GetChrome(id)
}

func GetInstance(id int) (*Chrome, error) {
	ins, err := instanceCache.Get(id)
	if err != nil {
		return nil, err
	} else {
		return ins.(*Chrome), nil
	}
}

func (i *Chrome) initialize() error {
	if i.State != STATE_UNINITIALIZED {
		return fmt.Errorf("instance in invalid state: %v", i.State)
	}

	ctx, cancel := context.WithTimeout(context.Background(), i.opts.InitTimeout)
	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(ctx, i.DebuggerURL)

	i.allocatorCtx = allocatorCtx
	i.cancelFunc = func() {
		allocatorCancel()
		cancel()
	}

	go i.stateManager()

	if ok, _ := RetryCheckChromeHealth(i.GetAddr(), 1, 0); !ok {
		i.cancelFunc()
		i.handleEvent(EVENT_INIT_FAIL)
		return fmt.Errorf("instance initialization failed: health check failed")
	}

	if err := i.handleEvent(EVENT_INIT_SUCCESS); err != nil {
		i.cancelFunc()
		return fmt.Errorf("failed to update state: %w", err)
	}

	go i.heartBeat()

	return nil
}

func (i *Chrome) GetNewContext() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(i.allocatorCtx)
}

func (i *Chrome) RetryInitialize(maxAttempts int) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := i.initialize(); err == nil {
			return nil
		} else {
			lastErr = err
			time.Sleep(time.Second * time.Duration(attempt)) // 指数退避
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

// 心跳检查
func (i *Chrome) heartBeat() {
	ticker := time.NewTicker(i.opts.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			state := i.GetState()
			if state != STATE_CONNECTED {
				continue
			}

			ok, _ := RetryCheckChromeHealth(i.GetAddr(), 1, 0)
			if !ok {
				i.handleEvent(EVENT_HEALTH_CHECK_FAIL)
			}

		case <-i.allocatorCtx.Done():
			return
		}
	}
}

func (i *Chrome) getAccounts() map[string]*account.Account {
	return i.accounts
}

func (i *Chrome) isAvailable(cat string) bool {
	account, exists := i.accounts[cat]
	if !exists {
		return false
	}

	return account.IsAvailable()
}

func (i *Chrome) GetAddr() string {
	return fmt.Sprintf("%s:%d", i.IP, i.Port)
}

func (i *Chrome) Close() error {
	if i.cancelFunc != nil {
		i.cancelFunc() // 取消 context，会导致 stateManager 和 heartBeat 退出
	}
	return nil
}

// 判断实例是否可用
func (i *Chrome) IsAvailable() bool {
	return i.GetState() == STATE_CONNECTED
}

// 状态管理器
func (i *Chrome) stateManager() {
	for {
		select {
		case evt := <-i.stateChan:
			var err error
			oldState := i.State

			switch evt.event {
			case EVENT_GET_STATE:
				evt.response <- i.State
				continue

			case EVENT_INIT_SUCCESS:
				if i.State == STATE_UNINITIALIZED {
					i.State = STATE_CONNECTED
				} else {
					err = fmt.Errorf("cannot initialize from state: %v", i.State)
				}

			case EVENT_INIT_FAIL:
				if i.State == STATE_UNINITIALIZED {
					i.State = STATE_INIT_FAILED
				} else {
					err = fmt.Errorf("cannot fail initialization from state: %v", i.State)
				}

			case EVENT_HEALTH_CHECK_SUCCESS:
				if i.State == STATE_DISCONNECTED {
					i.State = STATE_CONNECTED
				}

			case EVENT_HEALTH_CHECK_FAIL:
				if i.State == STATE_CONNECTED {
					i.State = STATE_DISCONNECTED
				}
			}

			if err == nil && oldState != i.State {
				// TODO: 更新实例状态
			}

			evt.response <- err

		case <-i.allocatorCtx.Done():
			return
		}
	}
}

// 处理状态事件
func (i *Chrome) handleEvent(event InstanceEvent) error {
	response := make(chan interface{}, 1)
	i.stateChan <- stateEvent{
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
func (i *Chrome) GetState() ChromeState {
	response := make(chan interface{}, 1)
	i.stateChan <- stateEvent{
		event:    EVENT_GET_STATE,
		response: response,
	}
	result := <-response
	return result.(ChromeState)
}

// NeedsReInitialize 判断是否需要重新初始化
func (i *Chrome) NeedsReInitialize() bool {
	state := i.GetState()
	return state == STATE_INIT_FAILED || state == STATE_DISCONNECTED
}

func (i *Chrome) cleanupTabs() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
	defer ticker.Stop()

	for {
		select {
		//case <-ticker.C:
		//	// 获取所有 targets (tabs)
		//	targets, err := chromedp.Targets(i.allocatorCtx)
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
		//				chromedp.CloseTarget(i.allocatorCtx, t.TargetID)
		//			}
		//		}
		//	}

		case <-i.allocatorCtx.Done():
			return
		}
	}
}

func (i *Chrome) ExecuteTask(task *task.Task) error {
	account, exists := i.accounts[task.Category]
	if !exists {
		return fmt.Errorf("no account found for category: %s", task.Category)
	}

	if !account.IsAvailable() {
		return fmt.Errorf("account not available")
	}

	select {
	case account.TaskChan <- task:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("send task timeout")
	}
}