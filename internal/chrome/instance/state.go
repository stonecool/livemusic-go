package instance

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome"
	"go.uber.org/zap"
)

// ChromeState 表示 Chrome 实例的状态
type ChromeState uint8

const (
	ChromeStateConnected    ChromeState = iota // 连接成功：包含初始化成功和心跳检查正常
	ChromeStateDisconnected                    // 连接断开：心跳检查失败
	ChromeStateOffline
)

// eventType 表示 Chrome 实例的事件类型
type eventType uint8

const (
	EventHealthCheckFail eventType = iota // 心跳检查失败
	EventGetState                         // 获取状态
)

// StateEvent 表示状态变更事件
type StateEvent struct {
	Type     eventType
	Response chan interface{}
}

// String 返回状态的字符串表示
func (s ChromeState) String() string {
	switch s {
	case ChromeStateConnected:
		return "Connected"
	case ChromeStateDisconnected:
		return "Disconnected"
	case ChromeStateOffline:
		return "Offline"
	default:
		return "Unknown"
	}
}

func (s ChromeState) IsValidTransition(event eventType) bool {
	switch event {
	case EventHealthCheckFail:
		return s == ChromeStateConnected
	default:
		return true
	}
}

func stateManager(c *Chrome) {
	for {
		select {
		case evt := <-c.StateChan:
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

func handleStateTransition(c *Chrome, evt StateEvent) {
	oldState := c.State
	var newState ChromeState
	var err error

	if !oldState.IsValidTransition(evt.Type) {
		err = fmt.Errorf("invalid state transition from %s with event %v",
			oldState.String(), evt.Type)
		evt.Response <- err
		return
	}

	switch evt.Type {
	case EventHealthCheckFail:
		newState = ChromeStateDisconnected
		//case eventShutdown:
		//	newState = stateOffline
	}

	if oldState != newState {
		c.State = newState

		if err := chrome.UpdateChrome(c); err != nil {
			internal.Logger.Error("failed to update chrome state",
				zap.Error(err),
				zap.Int("chromeID", c.ID),
				zap.String("oldState", oldState.String()),
				zap.String("newState", newState.String()))
		}
	}

	evt.Response <- err
}

func handleEvent(c *Chrome, event eventType) error {
	response := make(chan interface{}, 1)
	c.StateChan <- StateEvent{
		Type:     event,
		Response: response,
	}
	result := <-response
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}
