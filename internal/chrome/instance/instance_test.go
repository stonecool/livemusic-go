package instance

import (
	"context"
	"testing"
	"time"

	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stretchr/testify/assert"
)

func TestInstance_GetID(t *testing.T) {
	instance := &Instance{ID: 123}
	assert.Equal(t, 123, instance.GetID())
}

func TestInstance_GetAddr(t *testing.T) {
	instance := &Instance{
		IP:   "127.0.0.1",
		Port: 9222,
	}
	assert.Equal(t, "127.0.0.1:9222", instance.GetAddr())
}

func TestInstance_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		state    types.ChromeState
		expected bool
	}{
		{
			name:     "connected state",
			state:    types.ChromeStateConnected,
			expected: true,
		},
		{
			name:     "disconnected state",
			state:    types.ChromeStateDisconnected,
			expected: false,
		},
		{
			name:     "offline state",
			state:    types.ChromeStateOffline,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				State: tt.state,
				Opts: &types.InstanceOptions{
					InitTimeout:       time.Second * 100,
					HeartbeatInterval: time.Second,
				},
			}

			instance.initialize()

			assert.Equal(t, tt.expected, instance.getState() == types.ChromeStateConnected)
		})
	}
}

func TestInstance_HandleStateTransition(t *testing.T) {
	tests := []struct {
		name          string
		initialState  types.ChromeState
		event         types.EventType
		expectedState types.ChromeState
		expectError   bool
	}{
		{
			name:          "valid transition: connected to disconnected",
			initialState:  types.ChromeStateConnected,
			event:         types.EventHealthCheckFail,
			expectedState: types.ChromeStateDisconnected,
			expectError:   false,
		},
		{
			name:          "invalid transition",
			initialState:  types.ChromeStateConnected,
			event:         types.EventShutdown,
			expectedState: types.ChromeStateConnected,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				State:     tt.initialState,
				StateChan: make(chan types.StateEvent, 1),
			}

			go func() {
				instance.stateManager()
			}()

			time.Sleep(10 * time.Millisecond)

			err := instance.HandleEvent(tt.event)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedState, instance.State)
			}
		})
	}
}

func TestInstance_RetryInitialize(t *testing.T) {
	tests := []struct {
		name        string
		maxAttempts int
		opts        *types.InstanceOptions
		expectError bool
	}{
		{
			name:        "successful initialization",
			maxAttempts: 3,
			opts: &types.InstanceOptions{
				InitTimeout:       time.Second,
				HeartbeatInterval: time.Second,
			},
			expectError: false,
		},
		{
			name:        "initialization timeout",
			maxAttempts: 2,
			opts: &types.InstanceOptions{
				InitTimeout:       time.Nanosecond,
				HeartbeatInterval: time.Second,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				IP:        "127.0.0.1",
				Port:      9222,
				Opts:      tt.opts,
				StateChan: make(chan types.StateEvent, 1),
			}

			err := instance.RetryInitialize(tt.maxAttempts)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, instance.allocatorCtx)
				assert.NotNil(t, instance.cancelFunc)
			}

			if instance.cancelFunc != nil {
				instance.cancelFunc()
			}
		})
	}
}

func TestInstance_Close(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	instance := &Instance{
		allocatorCtx: ctx,
		cancelFunc:   cancel,
	}

	err := instance.Close()
	assert.NoError(t, err)

	select {
	case <-instance.allocatorCtx.Done():
	default:
		t.Error("Context was not cancelled")
	}
}

func TestInstance_GetNewContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	instance := &Instance{
		allocatorCtx: ctx,
	}

	newCtx, newCancel := instance.GetNewContext()
	defer newCancel()

	assert.NotNil(t, newCtx)
	assert.NotNil(t, newCancel)

	cancel()
	select {
	case <-newCtx.Done():
	default:
		t.Error("Child context was not cancelled with parent")
	}
}

func TestInstance_SetAccount(t *testing.T) {
	instance := &Instance{
		Accounts: make(map[string]account.IAccount),
	}

	mockAccount := &account.MockAccount{
		Category: "test",
	}

	instance.SetAccount("test", mockAccount)

	assert.Equal(t, mockAccount, instance.Accounts["test"])
}

func TestInstance_GetAccounts(t *testing.T) {
	instance := &Instance{
		Accounts: map[string]account.IAccount{
			"test1": &account.MockAccount{Category: "test1"},
			"test2": &account.MockAccount{Category: "test2"},
		},
	}

	accounts := instance.GetAccounts()
	assert.Len(t, accounts, 2)
	assert.Contains(t, accounts, "test1")
	assert.Contains(t, accounts, "test2")
}
