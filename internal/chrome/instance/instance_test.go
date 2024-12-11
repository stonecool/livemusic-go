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
			name:          "invalid to available not allowed",
			initialState:  types.InstanceStateInvalid,
			event:         types.EventHealthCheckSuccess,
			expectedState: types.InstanceStateInvalid,
			expectError:   true,
		},
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
			expectedState: types.InstanceStateAvailable,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				ID:        1,
				State:     tt.initialState,
				StateChan: make(chan types.StateEvent, 1),
			}

			// 发送事件
			evt := types.StateEvent{
				Type:     tt.event,
				Data:     tt.eventData,
				Response: make(chan interface{}, 1),
			}

			// 直接调用 HandleStateTransition
			instance.HandleStateTransition(evt)

			// 等待响应
			select {
			case result := <-evt.Response:
				if tt.expectError {
					// 检查是否是 error 类型
					if err, ok := result.(error); ok {
						assert.Error(t, err)
					} else {
						t.Error("expected error response, got:", result)
					}
				} else {
					// 如果期望成功，result 可能是 nil 或其他非错误值
					if err, ok := result.(error); ok {
						assert.NoError(t, err)
					}
					assert.Equal(t, tt.expectedState, instance.GetState())
				}
			case <-time.After(time.Second):
				t.Error("test timed out")
			}
		})
	}
}

func TestInstance_Initialize(t *testing.T) {
	instance := &Instance{
		ID:        1,
		StateChan: make(chan types.StateEvent),
		Opts: &types.InstanceOptions{
			InitTimeout:       time.Second,
			HeartbeatInterval: time.Second,
		},
	}

	instance.Initialize()
	assert.NotNil(t, instance.allocatorCtx)
	assert.NotNil(t, instance.cancelFunc)

	// Cleanup
	instance.Close()
}

func TestInstance_Close(t *testing.T) {
	instance := &Instance{
		ID:        1,
		StateChan: make(chan types.StateEvent),
	}

	// Test close without initialization
	instance.Close() // Should not panic

	// Test close after initialization
	ctx, cancel := context.WithCancel(context.Background())
	instance.allocatorCtx = ctx
	instance.cancelFunc = cancel

	instance.Close()
	select {
	case <-instance.allocatorCtx.Done():
		// Success
	default:
		t.Error("context was not cancelled")
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
