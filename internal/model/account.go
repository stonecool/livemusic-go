package model

import (
	"gorm.io/gorm"
)

type Account struct {
	Model

	AccountType string
	AccountId   string
	AccountName string
	Cookies     string
	State       uint8
}

// AddAccount Adds a new crawl account
func AddAccount(data map[string]interface{}) (*Account, error) {
	account := Account{
		AccountType: data["account_type"].(string),
		State:       data["state"].(uint8),
	}

	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccount Gets a crawl account by id
func GetAccount(id int) (*Account, error) {
	var account Account
	if err := db.Where("id = ?", id).First(&account).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &account, nil
}

// GetCrawlAccounts Gets all account
func GetAccountsByType(accountType string) ([]*Account, error) {
	var accounts []Account
	if err := db.Where("deleted_at != 0 AND account_type = ?", accountType).Find(&accounts).Error; err != nil {
		return nil, err
	}

	ret := make([]*Account, len(accounts), len(accounts))
	for i, _ := range accounts {
		ret[i] = &accounts[i]
	}

	return ret, nil
}
