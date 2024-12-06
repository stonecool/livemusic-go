package types

// EventType 事件类型
type EventType uint8

const (
	EventUnknown EventType = iota
	EventHealthCheckFail
)

// String 返回事件类型的字符串表示
func (e EventType) String() string {
	switch e {
	case EventHealthCheckFail:
		return "HealthCheckFail"
	default:
		return "Unknown"
	}
}

// StateEvent 状态事件
type StateEvent struct {
	Type     EventType
	Response chan interface{}
	Data     interface{} // 可选的事件数据
}
