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
		message.AccountState_New: {
			message.AccountCmd_Invalid: message.AccountState_NotLoggedIn,
		},
		message.AccountState_NotLoggedIn: {
			message.AccountCmd_Login: message.AccountState_Ready,
		},
		message.AccountState_Ready: {
			message.AccountCmd_Crawl: message.AccountState_Running,
		},
		message.AccountState_Running: {
			message.AccountCmd_CrawlAck: message.AccountState_Ready,
		},
		message.AccountState_Expired: {
			message.AccountCmd_Login: message.AccountState_NotLoggedIn,
		},
	},
	ErrorTransitions: map[message.AccountState]message.AccountState{
		message.AccountState_NotLoggedIn: message.AccountState_New,
		message.AccountState_Ready:       message.AccountState_NotLoggedIn,
		message.AccountState_Running:     message.AccountState_Expired,
	},
}

var noLoginTransitions = transitions{
	CmdTransitions: map[message.AccountState]map[message.AccountCmd]message.AccountState{
		message.AccountState_New: {
			message.AccountCmd_Invalid: message.AccountState_Ready,
		},
		message.AccountState_Ready: {
			message.AccountCmd_Crawl: message.AccountState_Running,
		},
		message.AccountState_Running: {
			message.AccountCmd_CrawlAck: message.AccountState_Ready,
		},
	},
	ErrorTransitions: map[message.AccountState]message.AccountState{
		message.AccountState_Ready:   message.AccountState_NotLoggedIn,
		message.AccountState_Running: message.AccountState_Expired,
	},
}
