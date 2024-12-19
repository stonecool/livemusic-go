package pool

import (
	"context"

	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/task"
	"github.com/stretchr/testify/mock"
)

type MockChrome struct {
	mock.Mock
}

func (m *MockChrome) GetID() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockChrome) GetAddr() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockChrome) Initialize() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockChrome) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockChrome) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockChrome) GetState() types.InstanceState {
	args := m.Called()
	return args.Get(0).(types.InstanceState)
}

func (m *MockChrome) SetState(state types.InstanceState) {
	m.Called(state)
}

func (m *MockChrome) GetStateChan() chan types.StateEvent {
	args := m.Called()
	return args.Get(0).(chan types.StateEvent)
}

func (m *MockChrome) ExecuteTask(task task.ITask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockChrome) GetNewContext() (context.Context, context.CancelFunc) {
	args := m.Called()
	return args.Get(0).(context.Context), args.Get(1).(context.CancelFunc)
}
