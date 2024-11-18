package account

import "github.com/stonecool/livemusic-go/internal/database"

func CreateAccount(category string) (IAccount, error) {
	repo := NewRepositoryDB(database.DB)
	factory := newFactory(repo)
	return factory.createAccount(category)
}

func getAccount(id int) (IAccount, error) {
	repo := NewRepositoryDB(database.DB)
	account, err := repo.Get(id)
	if err != nil {
		return nil, err
	}

	switch account.Category {
	case "wechat":
		wechatAccount := &WeChatAccount{Account: account}
		wechatAccount.Init()

		return wechatAccount, nil
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
