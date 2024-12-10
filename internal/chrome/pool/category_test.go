package pool

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stretchr/testify/assert"
)

func TestCategory_AddChrome(t *testing.T) {
	cat := newCategory("test-category")

	mockChrome := new(MockChrome)
	mockChrome.On("GetID").Return(1)
	mockChrome.On("GetAddr").Return("127.0.0.1:9222")
	mockChrome.On("Initialize").Return(nil)
	mockChrome.On("Close").Return(nil)
	mockChrome.On("IsAvailable").Return(true)
	mockChrome.On("GetState").Return(types.ChromeStateConnected)
	mockChrome.On("getStateChan").Return(make(chan types.StateEvent))

	// Test adding new chrome
	cat.AddChrome(mockChrome)
	assert.Len(t, cat.chromes, 1)

	// Test adding duplicate chrome
	cat.AddChrome(mockChrome)
	assert.Len(t, cat.chromes, 1)
}

func TestCategory_GetChromes(t *testing.T) {
	cat := newCategory("test-category")

	// Test empty category
	chromes := cat.GetChromes()
	assert.Empty(t, chromes)

	// Add chrome and test
	mockChrome := new(MockChrome)
	mockChrome.On("GetID").Return(1)
	mockChrome.On("GetAddr").Return("127.0.0.1:9222")

	cat.AddChrome(mockChrome)
	chromes = cat.GetChromes()
	assert.Len(t, chromes, 1)
	assert.Equal(t, mockChrome, chromes[0])
}

func TestCategory_ContainChrome(t *testing.T) {
	cat := newCategory("test-category")

	mockChrome := new(MockChrome)
	mockChrome.On("GetID").Return(1)
	mockChrome.On("GetAddr").Return("127.0.0.1:9222")

	// Test non-existent chrome
	assert.False(t, cat.ContainChrome("127.0.0.1:9222"))

	// Add chrome and test
	cat.AddChrome(mockChrome)
	assert.True(t, cat.ContainChrome("127.0.0.1:9222"))
}
