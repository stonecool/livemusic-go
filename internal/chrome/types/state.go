package types

// InstanceState represents the state of a Chrome instance
type InstanceState uint8

const (
	InstanceStateAvailable InstanceState = iota
	InstanceStateUnstable
	InstanceStateUnavailable
)

// String returns the string representation of the state
func (s InstanceState) String() string {
	switch s {
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
