package chrome

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal"
	"go.uber.org/zap"
)

// chromeState 表示 Chrome 实例的状态
type chromeState uint8

const (
	chromeStateConnected    chromeState = iota // 连接成功：包含初始化成功和心跳检查正常
	chromeStateDisconnected                    // 连接断开：心跳检查失败
	chromeStateOffline
)

// eventType 表示 Chrome 实例的事件类型
type eventType uint8

const (
	EventHealthCheckFail eventType = iota // 心跳检查失败
	EventGetState                         // 获取状态
)

// stateEvent 表示状态变更事件
type stateEvent struct {
	Type     eventType
	Response chan interface{}
}

// String 返回状态的字符串表示
func (s chromeState) String() string {
	switch s {
	case chromeStateConnected:
		return "Connected"
	case chromeStateDisconnected:
		return "Disconnected"
	case chromeStateOffline:
		return "Offline"
	default:
		return "Unknown"
	}
}

func (s chromeState) IsValidTransition(event eventType) bool {
	switch event {
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

func handleEvent(c *Chrome, event eventType) error {
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
