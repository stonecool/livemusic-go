package instance

import (
	"context"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestInstance_HandleStateTransition(t *testing.T) {
	tests := []struct {
		name          string
		initialState  types.InstanceState
		event         types.EventType
		eventData     interface{}
		expectedState types.InstanceState
		expectError   bool
	}{
		{
			name:          "available to unstable on health check fail",
			initialState:  types.InstanceStateAvailable,
			event:         types.EventHealthCheckFail,
			eventData:     1, // first failure
			expectedState: types.InstanceStateUnstable,
			expectError:   false,
		},
		{
			name:          "unstable to unavailable on third health check fail",
			initialState:  types.InstanceStateUnstable,
			event:         types.EventHealthCheckFail,
			eventData:     3, // third failure
			expectedState: types.InstanceStateUnavailable,
			expectError:   false,
		},
		{
			name:          "unstable to available on health check success",
			initialState:  types.InstanceStateUnstable,
			event:         types.EventHealthCheckSuccess,
			eventData:     nil,
			expectedState: types.InstanceStateAvailable,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				State:     tt.initialState,
				StateChan: make(chan types.StateEvent, 1),
			}

			responseChan := make(chan interface{}, 1)
			instance.HandleStateTransition(types.StateEvent{
				Type:     tt.event,
				Data:     tt.eventData,
				Response: responseChan,
			})

			result := <-responseChan
			if tt.expectError {
				assert.NotNil(t, result)
				assert.Error(t, result.(error))
			} else {
				// 如果不期望错误，result 可能是 nil
				if result != nil {
					assert.NoError(t, result.(error))
				}
				assert.Equal(t, tt.expectedState, instance.GetState())
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
