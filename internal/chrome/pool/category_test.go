package pool

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stretchr/testify/assert"
)

func TestCategory_AddChrome(t *testing.T) {
	cat := newCategory("test-category")

	chrome := &mockChrome{
		addr: "127.0.0.1:9222",
		accounts: map[string]account.IAccount{
			"test-category": &account.MockAccount{Category: "test-category"},
		},
		state: types.InstanceStateAvailable,
	}

	// Test adding new chrome
	cat.AddChrome(chrome)
	assert.Len(t, cat.chromes, 1)

	// Test adding duplicate chrome
	cat.AddChrome(chrome)
	assert.Len(t, cat.chromes, 1)
}

func TestCategory_GetChromes(t *testing.T) {
	cat := newCategory("test-category")

	// Test empty category
	chromes := cat.GetChromes()
	assert.Empty(t, chromes)

	// Add chrome and test
	chrome := &mockChrome{
		addr: "127.0.0.1:9222",
		accounts: map[string]account.IAccount{
			"test-category": &account.MockAccount{Category: "test-category"},
		},
		state: types.InstanceStateAvailable,
	}

	cat.AddChrome(chrome)
	chromes = cat.GetChromes()
	assert.Len(t, chromes, 1)
	assert.Equal(t, chrome, chromes[0])
}

func TestCategory_ContainChrome(t *testing.T) {
	cat := newCategory("test-category")

	chrome := &mockChrome{
		addr: "127.0.0.1:9222",
		accounts: map[string]account.IAccount{
			"test-category": &account.MockAccount{Category: "test-category"},
		},
		state: types.InstanceStateAvailable,
	}

	// Test non-existent chrome
	assert.False(t, cat.ContainChrome("127.0.0.1:9222"))

	// Add chrome and test
	cat.AddChrome(chrome)
	assert.True(t, cat.ContainChrome("127.0.0.1:9222"))
}
