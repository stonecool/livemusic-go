package internal

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
)

type CrawlAccount struct {
	ID          int    `json:"id"`
	AccountType string `json:"account_type"`
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	cookies     []byte
}

func (a *CrawlAccount) init(m *model.CrawlAccount) {
	a.ID = m.ID
	a.AccountType = m.AccountType
	a.AccountId = m.AccountId
	a.AccountName = m.AccountName
	a.cookies = m.Cookies
}

func (a *CrawlAccount) Add() error {
	_, ok := config.AccountMap[a.AccountType]
	if !ok {
		return fmt.Errorf("account_type:%s not exists", a.AccountType)
	}

	data := map[string]interface{}{
		"account_type": a.AccountType,
	}

	if account, err := model.AddCrawlAccount(data); err != nil {
		return err
	} else {
		a.init(account)
		return nil
	}
}

func (a *CrawlAccount) Get() error {
	if account, err := model.GetCrawlAccount(a.ID); err != nil {
		return err
	} else {
		a.init(account)
		return nil
	}
}

func (a *CrawlAccount) GetAll() ([]*CrawlAccount, error) {
	if accounts, err := model.GetCrawlAccountAll(); err != nil {
		return nil, err
	} else {
		var s []*CrawlAccount

		for _, account := range accounts {
			tempAccount := &CrawlAccount{}
			tempAccount.init(account)
			s = append(s, tempAccount)
		}

		return s, nil
	}
}

func (a *CrawlAccount) Delete() error {
	account, err := model.GetCrawlAccount(a.ID)
	if err != nil {
		return err
	}

	return model.DeleteCrawlAccount(account)
}
