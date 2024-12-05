package storage

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CRUD(t *testing.T) {
	repo := NewMockRepository()

	// Test Create
	model, err := repo.Create("127.0.0.1", 9222, "ws://127.0.0.1:9222", types.ChromeStateConnected)
	assert.NoError(t, err)
	assert.NotNil(t, model)
	assert.Equal(t, 1, model.ID)

	// Test Get
	retrieved, err := repo.Get(model.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.IP, retrieved.IP)
	assert.Equal(t, model.Port, retrieved.Port)
	assert.Equal(t, model.DebuggerURL, retrieved.DebuggerURL)
	assert.Equal(t, model.State, retrieved.State)

	// Test Update
	model.State = int(types.ChromeStateDisconnected)
	err = repo.Update(model)
	assert.NoError(t, err)

	updated, err := repo.Get(model.ID)
	assert.NoError(t, err)
	assert.Equal(t, int(types.ChromeStateDisconnected), updated.State)

	// Test GetAll
	models, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, models, 1)

	// Test ExistsByIPAndPort
	exists, err := repo.ExistsByIPAndPort("127.0.0.1", 9222)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByIPAndPort("127.0.0.1", 9223)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test Delete
	err = repo.Delete(model.ID)
	assert.NoError(t, err)

	_, err = repo.Get(model.ID)
	assert.Error(t, err)
}

func TestRepository_CreateInvalid(t *testing.T) {
	repo := NewMockRepository()

	// Test invalid IP
	_, err := repo.Create("", 9222, "ws://127.0.0.1:9222", types.ChromeStateConnected)
	assert.Error(t, err)

	// Test invalid port
	_, err = repo.Create("127.0.0.1", 0, "ws://127.0.0.1:9222", types.ChromeStateConnected)
	assert.Error(t, err)

	// Test invalid debugger URL
	_, err = repo.Create("127.0.0.1", 9222, "", types.ChromeStateConnected)
	assert.Error(t, err)
}

func TestRepository_UpdateInvalid(t *testing.T) {
	repo := NewMockRepository()

	// Create a valid model first
	model, err := repo.Create("127.0.0.1", 9222, "ws://127.0.0.1:9222", types.ChromeStateConnected)
	assert.NoError(t, err)

	// Test update with invalid IP
	model.IP = ""
	err = repo.Update(model)
	assert.Error(t, err)

	// Test update with invalid port
	model.IP = "127.0.0.1"
	model.Port = 0
	err = repo.Update(model)
	assert.Error(t, err)

	// Test update non-existent model
	model.ID = 999
	err = repo.Update(model)
	assert.Error(t, err)
}
