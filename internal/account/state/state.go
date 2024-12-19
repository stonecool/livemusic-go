package state

import "github.com/stonecool/livemusic-go/internal/message"

type transitions struct {
	ValidTransitions map[message.AccountState][]message.AccountState
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

func (m *Manager) IsValidTransition(from, to message.AccountState) bool {
	validStates, exists := m.transitions.ValidTransitions[from]
	if !exists {
		return false
	}
	for _, validState := range validStates {
		if to == validState {
			return true
		}
	}
	return false
}

var defaultTransitions = transitions{
	ValidTransitions: map[message.AccountState][]message.AccountState{
		message.AccountState_Undefined:   {},
		message.AccountState_New:         {message.AccountState_NotLoggedIn},
		message.AccountState_NotLoggedIn: {message.AccountState_Ready},
		message.AccountState_Ready:       {message.AccountState_Running, message.AccountState_Expired},
		message.AccountState_Running:     {message.AccountState_Ready, message.AccountState_Expired},
		message.AccountState_Expired:     {message.AccountState_NotLoggedIn},
	},
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
	ValidTransitions: map[message.AccountState][]message.AccountState{
		message.AccountState_Undefined:   {},
		message.AccountState_New:         {message.AccountState_Ready},
		message.AccountState_NotLoggedIn: {},
		message.AccountState_Ready:       {message.AccountState_Running, message.AccountState_Expired},
		message.AccountState_Running:     {message.AccountState_Ready, message.AccountState_Expired},
		message.AccountState_Expired:     {message.AccountState_Ready},
	},
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
