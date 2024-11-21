package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) get(id int) (*Account, error) {
	args := m.Called(id)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockRepository) create(category string) (*Account, error) {
	args := m.Called(category)
	return args.Get(0).(*Account), args.Error(1)
}

func TestCreateInstance(t *testing.T) {
	mockRepo := new(MockRepository)
	accountRepo = mockRepo

	expectedAccount := &Account{
		ID:       1,
		Category: "wechat",
	}

	mockRepo.On("create", "wechat").Return(expectedAccount, nil)
	mockRepo.On("get", 1).Return(expectedAccount, nil)

	instance, err := CreateInstance("wechat")

	assert.NoError(t, err)
	assert.Equal(t, expectedAccount.ID, instance.GetID())
	assert.Equal(t, expectedAccount.Category, "wechat")
}
