package crawl

import (
	"github.com/gocolly/colly"
	"github.com/stonecool/livemusic-go/internal/model"
	"log"
	"reflect"
)

type Account struct {
	ID          int                    `json:"id"`
	AccountType string                 `json:"account_type"`
	AccountId   string                 `json:"account_id"`
	AccountName string                 `json:"account_name"`
	Headers     map[string]interface{} `json:"headers"`
	State       uint8                  `json:"state"`
}

func (a *Account) Add() error {
	account := map[string]interface{}{
		"account_type": a.AccountType,
		"state":        uint8(0),
	}

	if templateId, err := model.AddCrawlAccount(account); err != nil {
		return err
	} else {
		a.ID = templateId
		return nil
	}
}

// GetHeaders
func (a *Account) GetHeaders() (map[string]interface{}, error) {
	return a.Headers, nil
}

// GetQueryParams
func (a *Account) GetQueryParams() (map[string]interface{}, error) {
	params := make(map[string]interface{})
	//if err := json.Unmarshal([]byte(a.QueryParams), &params); err != nil {
	//	return nil, err
	//}

	return params, nil
}

// GetFormData
func (a *Account) GetFormData() string {
	return ""
}

func GetCrawlAccountByID(id int) (interface{}, error) {
	var p *model.CrawlAccount
	p, err := model.GetCrawlAccount(id)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	if reflect.ValueOf(*p).IsZero() {
		return &Account{}, nil
	}

	account := Account{
		ID:      p.ID,
		Headers: make(map[string]interface{}),
	}

	return &account, nil
}

func (a *Account) GetId() int {
	return a.ID
}

func (a *Account) GetChan() chan CmdRequest {
	return nil
}

func (a *Account) GetState() State {
	return State(a.State)
}

func (a *Account) SetState(state State) {
	a.State = uint8(state)
}

func (a *Account) Login() (bool, error) {
	return false, nil
}

func (a *Account) Crawl(instance *Instance) error {
	return nil
}

func (a *Account) GetLoginRequestData() []byte {
	return nil
}

func (a *Account) LoginRequestCallback(request *colly.Request) error {
	return nil
}

func (a *Account) LoginResponseCallback(response *colly.Response) error {
	return nil
}
