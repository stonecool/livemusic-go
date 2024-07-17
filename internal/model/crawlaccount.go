package model

type CrawlAccount struct {
	Model

	CrawlType   string // 微信公众号，微博
	AccountId   string
	AccountName string
	Cookies     []byte
}

// AddCrawlAccount Adds a new crawl account
func AddCrawlAccount(accountType string) (*CrawlAccount, error) {
	account := CrawlAccount{
		CrawlType: accountType,
	}

	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

// DeleteCrawlAccount Deletes a crawl account
func DeleteCrawlAccount(account *CrawlAccount) error {
	return db.Delete(account).Error
}

// GetCrawlAccount Gets a crawl account
func GetCrawlAccount(id int) (*CrawlAccount, error) {
	var account CrawlAccount
	if err := db.Where("id = ? AND deleted_at != 0", id).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
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
