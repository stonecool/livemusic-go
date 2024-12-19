package state

import (
	"fmt"
	"sync"

	"github.com/stonecool/livemusic-go/internal/message"
)

type Handler interface {
	GetState() message.AccountState
	SetState(message.AccountState)
	HandleStateTransition(cmd message.AccountCmd) error
	HandleError(err error)
}

func NewStateHandler(category string) *AccountStateHandler {
	return &AccountStateHandler{
		state: message.AccountState_New,
		mgr:   SelectStateManager(category),
	}
}

func SelectStateManager(category string) *Manager {
	switch category {
	case "noLogin":
		return &Manager{transitions: noLoginTransitions}
	default:
		return &Manager{transitions: defaultTransitions}
	}
}

type AccountStateHandler struct {
	state message.AccountState
	mu    sync.RWMutex
	mgr   *Manager
}

func (h *AccountStateHandler) GetState() message.AccountState {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.state
}

func (h *AccountStateHandler) SetState(state message.AccountState) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.state = state
}

func (h *AccountStateHandler) HandleStateTransition(cmd message.AccountCmd) error {
	currentState := h.GetState()

	nextState := h.mgr.GetNextState(currentState, cmd)

	if !h.mgr.IsValidTransition(currentState, nextState) {
		return fmt.Errorf("invalid state transition from %v to %v for command %v",
			currentState, nextState, cmd)
	}

	h.SetState(nextState)
	return nil
}

func (h *AccountStateHandler) HandleError(err error) {
	currentState := h.GetState()
	errorState := h.mgr.GetErrorState(currentState)

	if h.mgr.IsValidTransition(currentState, errorState) {
		h.SetState(errorState)
	}
}
