package model

import "time"

type CrawlAccount struct {
	Model

	Category    string
	AccountName string
	LastURL     string
	Cookies     []byte
	InstanceID  int
	State       int
}

// AddCrawlAccount Adds a new crawl account
func AddCrawlAccount(data map[string]interface{}) (*CrawlAccount, error) {
	account := CrawlAccount{
		Category:    data["category"].(string),
		AccountName: data["account_name"].(string),
		LastURL:     data["last_url"].(string),
		Cookies:     data["cookies"].([]byte),
		InstanceID:  data["instance_id"].(int),
		State:       data["status"].(int),
	}

	if err := DB.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

// GetCrawlAccount Gets a crawl account
func GetCrawlAccount(id int) (*CrawlAccount, error) {
	var account CrawlAccount
	// FIXME
	if err := DB.Where("id = ? AND deleted_at = ?", id, 0).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func GetCrawlAccountAll() ([]*CrawlAccount, error) {
	var accounts []*CrawlAccount
	if err := DB.Where("deleted_at = ?", 0).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func EditCrawlAccount(id int, data interface{}) error {
	return DB.Model(&CrawlAccount{}).Where("id = ? AND deleted_at = ?", id, 0).Updates(data).Error
}

// DeleteCrawlAccount Deletes a crawl account
func DeleteCrawlAccount(account *CrawlAccount) error {
	return DB.Model(account).Where("deleted_at = ?", 0).Update("deleted_at", time.Now().Unix()).Error
}

// GetCrawlAccountsByType Gets crawl accounts by type
func GetCrawlAccountsByType(crawlType string) ([]*CrawlAccount, error) {
	var crawls []CrawlAccount
	if err := DB.Where("deleted_at = 0 AND account_type = ?", crawlType).Find(&crawls).Error; err != nil {
		return nil, err
	}

	ret := make([]*CrawlAccount, len(crawls), len(crawls))
	for i, _ := range crawls {
		ret[i] = &crawls[i]
	}

	return ret, nil
}
