package types

// Chrome 实例状态
type ChromeState uint8

const (
	ChromeStateConnected ChromeState = iota
	ChromeStateDisconnected
	ChromeStateOffline
)

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

// 建议添加状态转换管理器接口
type StateManager interface {
	GetState() ChromeState
	SetState(state ChromeState)
	GetStateChan() chan StateEvent
	HandleStateTransition(evt StateEvent)
	HandleEvent(event EventType) error
}

// 添加状态转换规则定义
var validTransitions = map[ChromeState][]EventType{
	ChromeStateConnected:    {EventHealthCheckFail},
	ChromeStateDisconnected: {EventGetState},
	ChromeStateOffline:      {EventGetState},
}

// 优化状态转换验证
func (s ChromeState) IsValidTransition(event EventType) bool {
	validEvents, exists := validTransitions[s]
	if !exists {
		return false
	}
	for _, e := range validEvents {
		if e == event {
			return true
		}
	}
	return false
}
