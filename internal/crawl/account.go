package crawl

import (
	"github.com/gocolly/colly"
	"github.com/stonecool/livemusic-go/internal/model"
	"log"
	"reflect"
)

type Account struct {
	ID          int    `json:"id"`
	AccountType string `json:"account_type"`
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	cookies     map[string]interface{}
	State       uint8 `json:"state"`
}

func AddAccount(accountType string) (*Account, error) {
	data := map[string]interface{}{
		"account_type": accountType,
		"state":        AccountStateUnLogged,
	}

	if m, err := model.AddAccount(data); err != nil {
		return nil, err
	} else {
		account := Account{
			ID:          m.ID,
			AccountType: m.AccountType,
			State:       m.State,
		}

		return &account, nil
	}
}

func GetAccountByID(id int) (*Account, error) {
	m, err := model.GetAccount(id)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	if reflect.ValueOf(*m).IsZero() {
		return &Account{}, nil
	}

	account := Account{
		ID: m.ID,
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

func (a *Account) GetCookies() []byte {
	return nil
}
