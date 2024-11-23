package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) get(id int) (*account, error) {
	args := m.Called(id)
	return args.Get(0).(*account), args.Error(1)
}

func (m *MockRepository) create(category string, state state) (*account, error) {
	args := m.Called(category, state)
	return args.Get(0).(*account), args.Error(1)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	accountRepo = mockRepo

	expectedAccount := &account{
		ID:       1,
		Category: "wechat",
	}

	mockRepo.On("create", "wechat", stateNew).Return(expectedAccount, nil)
	mockRepo.On("get", 1).Return(expectedAccount, nil)

	instance, err := CreateAccount("wechat")

	assert.NoError(t, err)
	assert.Equal(t, expectedAccount.ID, instance.GetID())
	assert.Equal(t, expectedAccount.Category, "wechat")
}
