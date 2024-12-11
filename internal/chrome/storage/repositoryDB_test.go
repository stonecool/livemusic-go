package storage

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Get(id int) (*types.Model, error) {
	args := m.Called(id)
	return args.Get(0).(*types.Model), args.Error(1)
}

func (m *MockDB) Create(model *types.Model) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockDB) Update(model *types.Model) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockDB) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) GetAll() ([]*types.Model, error) {
	args := m.Called()
	return args.Get(0).([]*types.Model), args.Error(1)
}

func (m *MockDB) FindBy(query string, args ...interface{}) ([]*types.Model, error) {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).([]*types.Model), callArgs.Error(1)
}

func (m *MockDB) ExistsBy(query string, args ...interface{}) (bool, error) {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(bool), callArgs.Error(1)
}

func TestRepositoryDB_Get(t *testing.T) {
	mockDB := new(MockDB)
	repo := &repositoryDB{db: mockDB}

	expectedModel := &types.Model{
		IP:          "127.0.0.1",
		Port:        9222,
		DebuggerURL: "ws://127.0.0.1:9222",
		State:       int(types.InstanceStateAvailable),
	}

	mockDB.On("Get", 1).Return(expectedModel, nil)

	model, err := repo.Get(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.IP, model.IP)
	assert.Equal(t, expectedModel.Port, model.Port)
	assert.Equal(t, expectedModel.DebuggerURL, model.DebuggerURL)
	assert.Equal(t, expectedModel.State, model.State)
}

func TestRepositoryDB_Create(t *testing.T) {
	mockDB := new(MockDB)
	repo := &repositoryDB{db: mockDB}

	mockDB.On("Create", mock.AnythingOfType("*types.Model")).Return(nil)

	model, err := repo.Create("127.0.0.1", 9222, "ws://127.0.0.1:9222", types.InstanceStateAvailable)

	assert.NoError(t, err)
	assert.NotNil(t, model)
	assert.Equal(t, "127.0.0.1", model.IP)
	assert.Equal(t, 9222, model.Port)
	assert.Equal(t, "ws://127.0.0.1:9222", model.DebuggerURL)
	assert.Equal(t, int(types.InstanceStateAvailable), model.State)
}

func TestRepositoryDB_Update(t *testing.T) {
	mockDB := new(MockDB)
	repo := &repositoryDB{db: mockDB}

	model := &types.Model{
		IP:          "127.0.0.1",
		Port:        9222,
		DebuggerURL: "ws://127.0.0.1:9222",
		State:       int(types.InstanceStateAvailable),
	}

	mockDB.On("Update", model).Return(nil)

	err := repo.Update(model)
	assert.NoError(t, err)
}

func TestRepositoryDB_GetAll(t *testing.T) {
	mockDB := new(MockDB)
	repo := &repositoryDB{db: mockDB}

	expectedModels := []*types.Model{
		{
			IP:          "127.0.0.1",
			Port:        9222,
			DebuggerURL: "ws://127.0.0.1:9222",
			State:       int(types.InstanceStateAvailable),
		},
	}

	mockDB.On("GetAll").Return(expectedModels, nil)

	models, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expectedModels, models)
}

func TestRepositoryDB_CreateInvalid(t *testing.T) {
	mockDB := new(MockDB)
	repo := &repositoryDB{db: mockDB}

	// Test invalid IP
	_, err := repo.Create("", 9222, "ws://127.0.0.1:9222", types.InstanceStateAvailable)
	assert.Error(t, err)

	// Test invalid port
	_, err = repo.Create("127.0.0.1", 0, "ws://127.0.0.1:9222", types.InstanceStateAvailable)
	assert.Error(t, err)

	// Test invalid debugger URL
	_, err = repo.Create("127.0.0.1", 9222, "", types.InstanceStateAvailable)
	assert.Error(t, err)
}
