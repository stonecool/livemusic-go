package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/database"
)

var accountFactory *factory
var accountCache *cache.Memo

func init() {
	accountCache = cache.New(func(id int) (interface{}, error) {
		return accountFactory.Get(id)
	})
	repo := NewRepositoryDB(database.DB)
	accountFactory = newAccountFactory(repo)
}

func CreateAccount(category string) (IAccount, error) {
	account, err := accountFactory.Create(category)
	if err != nil {
		return nil, err
	}

	if err := accountCache.Set(account.GetId(), account); err != nil {
		return nil, fmt.Errorf("failed to cache account: %w", err)
	}

	return account, nil
}

func GetAccount(ID int) (IAccount, error) {
	if acc, err := accountCache.Get(ID); err != nil {
		return nil, err
	} else {
		return acc.(IAccount), nil
	}
}

func UpdateAccount(account *Account) error {
	repo := NewRepositoryDB(database.DB)
	return repo.Update(account)
}

func DeleteAccount(id int) error {
	repo := NewRepositoryDB(database.DB)
	return repo.Delete(id)
}
