package state

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account/types"
	"github.com/stonecool/livemusic-go/internal/message"
	"go.uber.org/zap"
	"sync"
)

func NewStateHandler(category string) *AccountStateHandler {
	return &AccountStateHandler{
		state: message.AccountState_AS_New,
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

func (h *AccountStateHandler) Transit(account types.Account, cmd message.AccountCmd) {
	currentState := h.GetState()
	nextState := h.mgr.GetNextState(currentState, cmd)

	internal.Logger.Info("Account state transition",
		zap.Int("account_id", account.GetID()),
		zap.String("command", cmd.String()),
		zap.String("from_state", currentState.String()),
		zap.String("to_state", nextState.String()),
	)

	h.SetState(nextState)
}

func (h *AccountStateHandler) HandleError(account types.Account, err error) {
	currentState := h.GetState()
	errorState := h.mgr.GetErrorState(currentState)

	internal.Logger.Error("Account error state transition",
		zap.Int("account_id", account.GetID()),
		zap.Error(err),
		zap.String("from_state", currentState.String()),
		zap.String("to_state", errorState.String()),
	)

	h.SetState(errorState)
}
