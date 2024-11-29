package types

import (
	"fmt"
	"time"
)

// Chrome 实例状态
type ChromeState uint8

const (
	ChromeStateConnected ChromeState = iota
	ChromeStateDisconnected
	ChromeStateOffline
)

// 实例配置选项
type InstanceOptions struct {
	HeartbeatInterval   time.Duration
	InitTimeout         time.Duration
	TabCleanupInterval  time.Duration
	TabInactiveTimeout  time.Duration
	ZombieCheckInterval time.Duration
}

// StateEvent 状态事件
type StateEvent struct {
	Type     EventType
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

func (s ChromeState) IsValidTransition(event EventType) bool {
	switch event {
	case EventHealthCheckFail:
		return s == ChromeStateConnected
	default:
		return true
	}
}

func HandleStateTransition(c IChrome, evt StateEvent) {
	oldState := c.GetState()
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

func HandleEvent(c IChrome, event EventType) error {
	response := make(chan interface{}, 1)
	c.GetStateChan() <- StateEvent{
		Type:     event,
		Response: response,
	}
	result := <-response
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

type EventType uint8

const (
	EventHealthCheckFail EventType = iota
	EventGetState
)
