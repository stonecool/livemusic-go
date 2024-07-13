package internal

import (
	"github.com/stonecool/livemusic-go/internal/model"
	"reflect"
)

type CrawlAccount struct {
	ID          int    `json:"id"`
	AccountType string `json:"account_type"`
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	State       uint8  `json:"state"`
	cookies     []byte
}

func initCrawlAccount(m *model.CrawlAccount) *CrawlAccount {
	if m == nil || reflect.ValueOf(m).IsZero() {
		return nil
	}

	var account CrawlAccount
	account.ID = m.ID
	account.AccountType = m.CrawlType
	account.AccountId = m.AccountId
	account.AccountName = m.AccountName
	account.cookies = m.Cookies

	return &account
}

// AddCrawlAccount
func AddCrawlAccount(accountType string) (*CrawlAccount, error) {
	_, ok := CrawlAccountMap[accountType]
	if !ok {
		return nil, error(nil)
	}

	data := map[string]interface{}{
		"account_type": accountType,
		"state":        uint8(0),
	}

	if m, err := model.AddCrawlAccount(data); err != nil {
		return nil, err
	} else {
		return initCrawlAccount(m), nil
	}
}

// GetCrawlAccount
func GetCrawlAccount(id int) (*CrawlAccount, error) {
	m, err := model.GetCrawlAccount(id)
	if err != nil {
		return nil, err
	} else {
		return initCrawlAccount(m), nil
	}
}
