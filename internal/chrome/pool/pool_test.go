package pool

import (
	"testing"

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
					"test-category": &account.MockAccount{Category: "test-category"},
				},
				state: types.InstanceStateAvailable,
			},
			wantErr: false,
		},
		{
			name: "Add duplicate chrome instance",
			chrome: &mockChrome{
				addr: "localhost:9222",
				accounts: map[string]account.IAccount{
					"test-category": &account.MockAccount{Category: "test-category"},
				},
				state: types.InstanceStateAvailable,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.AddChrome(tt.chrome)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "exists in pool")
			} else {
				assert.NoError(t, err)
				// 验证实例是否被添加到池中
				chrome := p.chromes[tt.chrome.GetAddr()]
				assert.NotNil(t, chrome)
				assert.Equal(t, tt.chrome.GetAddr(), chrome.GetAddr())
			}
		})
	}
}

func TestPool_GetChromesByCategory(t *testing.T) {
	p := GetPool()

	// Test empty category
	chromes := p.GetChromesByCategory("non-existent")
	assert.Empty(t, chromes)

	// Add chrome instance
	chrome := &mockChrome{
		addr: "127.0.0.1:9222",
		accounts: map[string]account.IAccount{
			"test-category": &account.MockAccount{Category: "test-category"},
		},
		state: types.InstanceStateAvailable,
	}

	// First add chrome to pool
	err := p.AddChrome(chrome)
	assert.NoError(t, err)

	// Then login to add it to category
	p.Login(chrome, "test-category")

	// Test getting chromes by category
	chromes = p.GetChromesByCategory("test-category")
	assert.Len(t, chromes, 1)
	assert.Equal(t, chrome, chromes[0])

	// Test getting chromes from another category
	chromes = p.GetChromesByCategory("non-existent-category")
	assert.Empty(t, chromes)
}

func TestPool_GetAllChromes(t *testing.T) {
	p := GetPool()

	// Test empty pool
	chromes := p.GetAllChromes()
	assert.Empty(t, chromes)

	// Add chrome instance
	chrome := &mockChrome{
		addr: "127.0.0.1:9222",
		accounts: map[string]account.IAccount{
			"test-category": &account.MockAccount{Category: "test-category"},
		},
		state: types.InstanceStateAvailable,
	}

	err := p.AddChrome(chrome)
	assert.NoError(t, err)

	// Test getting all chrome instances
	chromes = p.GetAllChromes()
	assert.Len(t, chromes, 1)
	assert.Contains(t, chromes, chrome.GetAddr())
	assert.Equal(t, chrome, chromes[chrome.GetAddr()])
}

func TestPool_Login(t *testing.T) {
	p := GetPool()

	// Add chrome instance
	mockChrome := &mockChrome{
		addr: "127.0.0.1:9222",
		accounts: map[string]account.IAccount{
			"test-category": &account.MockAccount{Category: "test-category"},
		},
		state: types.InstanceStateAvailable,
	}

	err := p.AddChrome(mockChrome)
	assert.NoError(t, err)

	// Test login
	p.Login(mockChrome, "test-category")
	chromes := p.GetChromesByCategory("test-category")
	assert.Len(t, chromes, 1)
	assert.Equal(t, mockChrome, chromes[0])
}
