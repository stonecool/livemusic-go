package internal

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
	"sync"
)

type CrawlAccount struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	AccountName  string `json:"account_name"`
	lastLoginURL string
	cookies      []byte
	instanceAddr string
	state        AccountState
	mu           sync.Mutex
}

func (ca *CrawlAccount) init(m *model.CrawlAccount) {
	ca.ID = m.ID
	ca.Category = m.Category
	ca.AccountName = m.AccountName
	ca.cookies = m.Cookies
	ca.lastLoginURL = m.LastLoginURL
}

func (ca *CrawlAccount) Add() error {
	_, ok := config.AccountMap[ca.Category]
	if !ok {
		return fmt.Errorf("account_type:%s not exists", ca.Category)
	}

	data := map[string]interface{}{
		"account_type": ca.Category,
	}

	if account, err := model.AddCrawlAccount(data); err != nil {
		return err
	} else {
		ca.init(account)
		return nil
	}
}

func (ca *CrawlAccount) Get() error {
	if account, err := model.GetCrawlAccount(ca.ID); err != nil {
		return err
	} else {
		ca.init(account)
		return nil
	}
}

func (ca *CrawlAccount) GetAll() ([]*CrawlAccount, error) {
	if accounts, err := model.GetCrawlAccountAll(); err != nil {
		return nil, err
	} else {
		var s []*CrawlAccount

		for _, account := range accounts {
			acc := &CrawlAccount{}
			acc.init(account)
			s = append(s, acc)
		}

		return s, nil
	}
}

func (ca *CrawlAccount) Edit() error {
	data := map[string]interface{}{
		"account_name":   ca.AccountName,
		"last_login_url": ca.lastLoginURL,
		"cookies":        ca.cookies,
	}

	return model.EditCrawlAccount(ca.ID, data)
}

func (ca *CrawlAccount) Delete() error {
	account, err := model.GetCrawlAccount(ca.ID)
	if err != nil {
		return err
	}

	return model.DeleteCrawlAccount(account)
}

func (ca *CrawlAccount) getCategory() string {
	return ca.Category
}

func (ca *CrawlAccount) IsAvailable() bool {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	return ca.state == AS_RUNNING
}
