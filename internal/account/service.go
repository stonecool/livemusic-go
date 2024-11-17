package account

import "github.com/stonecool/livemusic-go/internal/database"

func CreateAccount(category string) (IAccount, error) {
	repo := NewRepositoryDB(database.DB)
	factory := newFactory(repo)
	account, err := factory.createAccount(category)

	if err != nil {
		return nil, err
	}

	switch account.GetCategory() {
	case "wechat":
		return &WeChatAccount{Account: *account}, nil
	default:
		return account, nil
	}
}

func GetAccount(id int) (IAccount, error) {
	repo := NewRepositoryDB(database.DB)
	account, err := repo.Get(id)
	if err != nil {
		return nil, err
	}

	switch account.Category {
	case "wechat":
		return &WeChatAccount{Account: *account}, nil
	default:
		return account, nil
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
