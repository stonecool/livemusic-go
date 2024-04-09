package model

import (
	"gorm.io/gorm"
)

type CrawlAccount struct {
	Model

	TemplateId  string
	Name        string
	Username    string
	Password    string
	Headers     string
	QueryParams string
	FormData    string
	State       uint8
}

// AddCrawlAccount Adds a new crawl account
func AddCrawlAccount(data map[string]interface{}) (int, error) {
	account := CrawlAccount{
		TemplateId: data["template_id"].(string),
		Name:       data["name"].(string),
		Username:   data["username"].(string),
		Password:   data["password"].(string),
		Headers:    data["headers"].(string),
		State:      data["state"].(uint8),
	}

	if err := db.Create(&account).Error; err != nil {
		return 0, err
	}

	return account.ID, nil
}

// GetCrawlAccount Gets a crawl account by id
func GetCrawlAccount(id int) (*CrawlAccount, error) {
	var account CrawlAccount
	if err := db.Where("id = ?", id).First(&account).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &account, nil
}

// GetCrawlAccounts Gets all account
func GetCrawlAccounts() ([]*CrawlAccount, error) {
	var accounts []CrawlAccount
	if err := db.Where("deleted_at != 0").Find(&accounts).Error; err != nil {
		return nil, err
	}

	ret := make([]*CrawlAccount, len(accounts), len(accounts))
	for i, _ := range accounts {
		ret[i] = &accounts[i]
	}

	return ret, nil
}
