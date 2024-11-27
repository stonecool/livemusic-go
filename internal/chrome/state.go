package chrome

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal"
	"go.uber.org/zap"
)

// chromeState 表示 Chrome 实例的状态
type chromeState uint8

const (
	chromeStateUninitialized chromeState = iota // 未初始化：实例刚创建
	chromeStateInitFailed                       // 初始化失败：初始化方法执行失败
	chromeStateConnected                        // 连接成功：包含初始化成功和心跳检查正常
	chromeStateDisconnected                     // 连接断开：心跳检查失败
)

// EventType 表示 Chrome 实例的事件类型
type EventType uint8

const (
	EventStart              EventType = iota // 开始初始化
	EventInitSuccess                         // 初始化成功
	EventInitFail                            // 初始化失败
	EventHealthCheckSuccess                  // 心跳检查成功
	EventHealthCheckFail                     // 心跳检查失败
	EventGetState                            // 获取状态
)

// stateEvent 表示状态变更事件
type stateEvent struct {
	Type     EventType
	Response chan interface{}
}

// String 返回状态的字符串表示
func (s chromeState) String() string {
	switch s {
	case chromeStateUninitialized:
		return "Uninitialized"
	case chromeStateInitFailed:
		return "InitFailed"
	case chromeStateConnected:
		return "Connected"
	case chromeStateDisconnected:
		return "Disconnected"
	default:
		return "Unknown"
	}
}

func (s chromeState) IsValidTransition(event EventType) bool {
	switch event {
	case EventInitSuccess:
		return s == chromeStateUninitialized
	case EventInitFail:
		return s == chromeStateUninitialized
	case EventHealthCheckSuccess:
		return s == chromeStateDisconnected
	case EventHealthCheckFail:
		return s == chromeStateConnected
	default:
		return true
	}
}

func stateManager(c *Chrome) {
	for {
		select {
		case evt := <-c.stateChan:
			switch evt.Type {
			case EventGetState:
				evt.Response <- c.State
				continue
			default:
				handleStateTransition(c, evt)
			}
		case <-c.allocatorCtx.Done():
			return
		}
	}
}

func handleStateTransition(c *Chrome, evt stateEvent) {
	oldState := c.State
	var err error

	if !oldState.IsValidTransition(evt.Type) {
		err = fmt.Errorf("invalid state transition from %s with event %v",
			oldState.String(), evt.Type)
		evt.Response <- err
		return
	}

	switch evt.Type {
	case EventInitSuccess:
		c.State = chromeStateConnected
	case EventInitFail:
		c.State = chromeStateInitFailed
	case EventHealthCheckSuccess:
		c.State = chromeStateConnected
	case EventHealthCheckFail:
		c.State = chromeStateDisconnected
	}

	if oldState != c.State {
		// TODO: 更新数据库中的状态
		if err := repo.update(c); err != nil {
			internal.Logger.Error("failed to update chrome state",
				zap.Error(err),
				zap.Int("chromeID", c.ID),
				zap.String("oldState", oldState.String()),
				zap.String("newState", c.State.String()))
		}
	}

	evt.Response <- err
}

func handleEvent(c *Chrome, event EventType) error {
	response := make(chan interface{}, 1)
	c.stateChan <- stateEvent{
		Type:     event,
		Response: response,
	}
	result := <-response
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}
