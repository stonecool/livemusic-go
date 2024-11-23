package account

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Get(id int) (*accountModel, error) {
	args := m.Called(id)
	return args.Get(0).(*accountModel), args.Error(1)
}

func (m *MockDB) Create(model *accountModel) error {
	args := m.Called(model)
	return args.Error(0)
}

//func TestRepositoryDB_Get(t *testing.T) {
//    mockDB := new(MockDB)
//    repo := &repositoryDB{db: mockDB}
//
//    expectedModel := &accountModel{
//        Category: "wechat",
//        Name: "test",
//    }
//
//    mockDB.On("Get", 1).Return(expectedModel, nil)
//
//    account, err := repo.get(1)
//
//    assert.NoError(t, err)
//    assert.Equal(t, expectedModel.Category, account.Category)
//    assert.Equal(t, expectedModel.Name, account.Name)
//}
//
//func TestRepositoryDB_Create(t *testing.T) {
//    mockDB := new(MockDB)
//    repo := &repositoryDB{db: mockDB}
//
//    mockDB.On("Create", mock.AnythingOfType("*account.accountModel")).Return(nil)
//
//    account, err := repo.create("wechat")
//
//    assert.NoError(t, err)
//    assert.Equal(t, "wechat", account.Category)
//}
