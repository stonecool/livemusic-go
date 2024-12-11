package pool

import (
	"context"

	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/task"
)

// Mock implementations for testing
type mockChrome struct {
	addr     string
	accounts map[string]account.IAccount
	state    types.InstanceState
}

func (m *mockChrome) GetAddr() string                          { return m.addr }
func (m *mockChrome) GetID() int                               { return 1 }
func (m *mockChrome) GetAccounts() map[string]account.IAccount { return m.accounts }
func (m *mockChrome) GetNewContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
func (m *mockChrome) ExecuteTask(task task.ITask) error  { return nil }
func (m *mockChrome) Initialize()                        {}
func (m *mockChrome) Close()                             {}
func (m *mockChrome) IsAvailable() bool                  { return m.state == types.InstanceStateAvailable }
func (m *mockChrome) GetState() types.InstanceState      { return m.state }
func (m *mockChrome) SetState(state types.InstanceState) { m.state = state }
