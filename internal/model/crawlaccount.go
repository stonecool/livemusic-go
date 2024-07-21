package model

type CrawlAccount struct {
	Model

	AccountType string
	AccountId   string
	AccountName string
	Cookies     []byte
}

// AddCrawlAccount Adds a new crawl account
func AddCrawlAccount(data map[string]interface{}) (*CrawlAccount, error) {
	account := CrawlAccount{
		AccountType: data["account_type"].(string),
	}

	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

// GetCrawlAccount Gets a crawl account
func GetCrawlAccount(id int) (*CrawlAccount, error) {
	var account CrawlAccount
	// FIXME
	if err := db.Where("id = ? AND deleted_at != ?", id, 0).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func GetCrawlAccountAll() ([]*CrawlAccount, error) {
	var accounts []*CrawlAccount
	if err := db.Where("deleted_at != ?", 0).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

// DeleteCrawlAccount Deletes a crawl account
func DeleteCrawlAccount(account *CrawlAccount) error {
	return db.Delete(account).Error
}

// GetCrawlAccountsByType Gets crawl accounts by type
func GetCrawlAccountsByType(crawlType string) ([]*CrawlAccount, error) {
	var crawls []CrawlAccount
	if err := db.Where("deleted_at != 0 AND account_type = ?", crawlType).Find(&crawls).Error; err != nil {
		return nil, err
	}

	ret := make([]*CrawlAccount, len(crawls), len(crawls))
	for i, _ := range crawls {
		ret[i] = &crawls[i]
	}

	return ret, nil
}
