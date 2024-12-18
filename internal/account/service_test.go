package account

import (
	"github.com/stonecool/livemusic-go/internal/account/state"
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

func (m *MockRepository) create(category string, state state.accountState) (*account, error) {
	args := m.Called(category, state)
	return args.Get(0).(*account), args.Error(1)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	repo = mockRepo

	expectedAccount := &account{
		ID:       1,
		Category: "wechat",
	}

	mockRepo.On("create", "wechat", state.stateNew).Return(expectedAccount, nil)
	mockRepo.On("get", 1).Return(expectedAccount, nil)

	instance, err := CreateAccount("wechat")

	assert.NoError(t, err)
	assert.Equal(t, expectedAccount.ID, instance.GetID())
	assert.Equal(t, expectedAccount.Category, "wechat")
}
