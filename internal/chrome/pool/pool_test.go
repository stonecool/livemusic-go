package pool

import (
	"context"
	"testing"

	"github.com/stonecool/livemusic-go/internal/task"

	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stretchr/testify/assert"
)

func TestPool_AddChrome(t *testing.T) {
	p := GetPool()

	tests := []struct {
		name    string
		chrome  types.Chrome
		wantErr bool
	}{
		{
			name: "Add new chrome instance successfully",
			chrome: &mockChrome{
				addr: "localhost:9222",
				accounts: map[string]account.IAccount{
					"music": &account.MockAccount{Category: "music"},
					"video": &account.MockAccount{Category: "video"},
				},
			},
			wantErr: false,
		},
		{
			name: "Add duplicate chrome instance",
			chrome: &mockChrome{
				addr:     "localhost:9222",
				accounts: map[string]account.IAccount{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.AddChrome(tt.chrome)
			if (err != nil) != tt.wantErr {
				t.Errorf("pool.AddChrome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify chrome was added to the pool
				if _, exists := p.chromes[tt.chrome.GetAddr()]; !exists {
					t.Error("Chrome instance was not added to pool.chromes")
				}

				// Verify categories were created and chrome was added to them
				for _, acc := range tt.chrome.GetAccounts() {
					cat := acc.GetCategory()
					if _, exists := p.categories[cat]; !exists {
						t.Errorf("Category %s was not created", cat)
					}
					if !p.categories[cat].ContainChrome(tt.chrome.GetAddr()) {
						t.Errorf("Chrome was not added to category %s", cat)
					}
				}
			}
		})
	}
}

func TestPool_GetChromesByCategory(t *testing.T) {
	p := GetPool()

	// Test empty category
	chromes := p.GetChromesByCategory("non-existent")
	assert.Nil(t, chromes)

	// Add chrome to category
	mockChrome := new(MockChrome)
	mockChrome.On("GetID").Return(1)
	mockChrome.On("GetAddr").Return("127.0.0.1:9222")
	mockChrome.On("GetAccounts").Return(map[string]account.IAccount{
		"test-category": &account.MockAccount{Category: "test-category"},
	})
	mockChrome.On("Initialize").Return(nil)
	mockChrome.On("Close").Return(nil)
	mockChrome.On("IsAvailable").Return(true)
	mockChrome.On("GetState").Return(types.ChromeStateConnected)
	mockChrome.On("GetStateChan").Return(make(chan types.StateEvent))

	err := p.AddChrome(mockChrome)
	assert.NoError(t, err)

	// Test getting chromes by category
	chromes = p.GetChromesByCategory("test-category")
	assert.Len(t, chromes, 1)
	assert.Equal(t, mockChrome, chromes[0])
}

// Mock implementations for testing
type mockChrome struct {
	addr      string
	accounts  map[string]account.IAccount
	state     types.InstanceState
	stateChan chan types.StateEvent
}

func (m *mockChrome) GetAddr() string                          { return m.addr }
func (m *mockChrome) GetID() int                               { return 1 }
func (m *mockChrome) GetAccounts() map[string]account.IAccount { return m.accounts }
func (m *mockChrome) GetNewContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
func (m *mockChrome) ExecuteTask(task task.ITask) error { return nil }
func (m *mockChrome) Initialize() error                 { return nil }
func (m *mockChrome) Close() error                      { return nil }
func (m *mockChrome) IsAvailable() bool                 { return true }
func (m *mockChrome) GetState() types.InstanceState {
	return types.ChromeStateConnected
}

func (m *mockChrome) SetState(state types.InstanceState) {
	m.state = state
}

func (m *mockChrome) GetStateChan() chan types.StateEvent {
	if m.stateChan == nil {
		m.stateChan = make(chan types.StateEvent)
	}
	return m.stateChan
}
