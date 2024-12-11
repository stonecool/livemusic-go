package types

// InstanceState represents the state of a Chrome instance
type InstanceState uint8

const (
	InstanceStateInvalid     InstanceState = iota // 非法状态，用作零值
	InstanceStateAvailable                        // 正常可用状态
	InstanceStateUnstable                         // 临时不可用状态（单次心跳失败）
	InstanceStateUnavailable                      // 不可用状态（多次心跳失败）
)

// String returns the string representation of the state
func (s InstanceState) String() string {
	switch s {
	case InstanceStateInvalid:
		return "Invalid"
	case InstanceStateAvailable:
		return "Available"
	case InstanceStateUnstable:
		return "Unstable"
	case InstanceStateUnavailable:
		return "Unavailable"
	default:
		return "Unknown"
	}
}

// StateManager interface for managing Chrome instance states
type StateManager interface {
	GetState() InstanceState
	SetState(state InstanceState)
	GetStateChan() chan StateEvent
	HandleStateTransition(evt StateEvent)
	HandleEvent(event EventType) error
}

// validTransitions defines valid state transitions based on events
var validTransitions = map[InstanceState][]EventType{
	InstanceStateAvailable: {
		EventHealthCheckSuccess,
		EventHealthCheckFail,
	},
	InstanceStateUnstable: {
		EventHealthCheckSuccess, // 恢复到可用状态
		EventHealthCheckFail,    // 累计失败次数增加
	},
	InstanceStateUnavailable: {
		EventHealthCheckSuccess, // 允许从不可用恢复到可用状态
	},
	// 非法状态不允许任何转换
	InstanceStateInvalid: {},
}

// IsValidTransition checks if the state transition is valid for the given event
func (s InstanceState) IsValidTransition(event EventType) bool {
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
