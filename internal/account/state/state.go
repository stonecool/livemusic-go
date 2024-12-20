package state

import "github.com/stonecool/livemusic-go/internal/message"

type transitions struct {
	CmdTransitions   map[message.AccountState]map[message.AccountCmd]message.AccountState
	ErrorTransitions map[message.AccountState]message.AccountState
}

type Manager struct {
	transitions transitions
}

func (m *Manager) GetNextState(state message.AccountState, cmd message.AccountCmd) message.AccountState {
	if cmdStates, exists := m.transitions.CmdTransitions[state]; exists {
		if nextState, hasCmd := cmdStates[cmd]; hasCmd {
			return nextState
		}
	}
	return state
}

func (m *Manager) GetErrorState(state message.AccountState) message.AccountState {
	if nextState, exists := m.transitions.ErrorTransitions[state]; exists {
		return nextState
	}
	return state
}

var defaultTransitions = transitions{
	CmdTransitions: map[message.AccountState]map[message.AccountCmd]message.AccountState{
		message.AccountState_AS_New: {
			message.AccountCmd_AC_INITIALIZE: message.AccountState_AS_NotLoggedIn,
		},
		message.AccountState_AS_NotLoggedIn: {
			message.AccountCmd_AC_Login: message.AccountState_AS_Ready,
		},
		message.AccountState_AS_Ready: {
			message.AccountCmd_AC_EXPIRED: message.AccountState_AS_Expired,
		},
		message.AccountState_AS_Expired: {
			message.AccountCmd_AC_Login: message.AccountState_AS_Ready,
		},
	},
	ErrorTransitions: map[message.AccountState]message.AccountState{
		message.AccountState_AS_Ready: message.AccountState_AS_Expired,
	},
}

var noLoginTransitions = transitions{
	CmdTransitions: map[message.AccountState]map[message.AccountCmd]message.AccountState{
		message.AccountState_AS_New: {
			message.AccountCmd_AC_INITIALIZE: message.AccountState_AS_Ready,
		},
	},
	ErrorTransitions: map[message.AccountState]message.AccountState{},
}
