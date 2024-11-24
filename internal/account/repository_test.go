package account

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Get(id int) (*model, error) {
	args := m.Called(id)
	return args.Get(0).(*model), args.Error(1)
}

func (m *MockDB) Create(model *model) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockDB) Update(model *model) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockDB) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) GetAll() ([]*model, error) {
	args := m.Called()
	return args.Get(0).([]*model), args.Error(1)
}

func (m *MockDB) FindBy(query string, args ...interface{}) ([]*model, error) {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).([]*model), callArgs.Error(1)
}

func (m *MockDB) ExistsBy(query string, args ...interface{}) (bool, error) {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(bool), callArgs.Error(1)
}

func TestRepositoryDB_Get(t *testing.T) {
	mockDB := new(MockDB)
	repo := &repositoryDB{db: mockDB}

	expectedModel := &model{
		Category: "wechat",
		Name:     "test",
	}

	mockDB.On("Get", 1).Return(expectedModel, nil)

	acc, err := repo.get(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Category, acc.Category)
	assert.Equal(t, expectedModel.Name, acc.Name)
}

func TestRepositoryDB_Create(t *testing.T) {
	mockDB := new(MockDB)
	repo := &repositoryDB{db: mockDB}

	mockDB.On("Create", mock.AnythingOfType("*account.model")).Return(nil)

	acc, err := repo.create("wechat", stateNew)

	assert.NoError(t, err)
	assert.Equal(t, "wechat", acc.Category)
	assert.Equal(t, stateNew, acc.State)
}
