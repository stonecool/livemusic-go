package account

type state int

const (
	stateNew state = iota
	stateInitialized
	stateNotLoggedIn
	stateReady
	stateRunning
	stateTerminated
)

func (s state) String() string {
	switch s {
	case stateNew:
		return "new"
	case stateInitialized:
		return "initialized"
	case stateNotLoggedIn:
		return "notLoggedIn"
	case stateReady:
		return "ready"
	case stateRunning:
		return "running"
	case stateTerminated:
		return "terminated"
	default:
		return "unknown"
	}
}

func (s state) isValidTransition(target state) bool {
	switch s {
	case stateNew:
		return target == stateInitialized
	case stateInitialized:
		return target == stateNotLoggedIn
	case stateNotLoggedIn:
		return target == stateReady || target == stateTerminated
	case stateReady:
		return target == stateRunning || target == stateNotLoggedIn || target == stateTerminated
	case stateRunning:
		return target == stateReady || target == stateTerminated
	case stateTerminated:
		return false // 终止状态不能转换到其他状态
	default:
		return false
	}
}
