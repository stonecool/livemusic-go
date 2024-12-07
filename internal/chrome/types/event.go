package types

// EventType represents the type of state transition event
type EventType uint8

const (
	EventUnknown            EventType = iota
	EventHealthCheckSuccess           // 心跳检查成功
	EventHealthCheckFail              // 心跳检查失败
)

// String returns the string representation of the event type
func (e EventType) String() string {
	switch e {
	case EventHealthCheckSuccess:
		return "HealthCheckSuccess"
	case EventHealthCheckFail:
		return "HealthCheckFail"
	default:
		return "Unknown"
	}
}

// StateEvent represents a state transition event with optional data
type StateEvent struct {
	Type     EventType
	Response chan interface{}
	Data     interface{} // Optional event data
}
