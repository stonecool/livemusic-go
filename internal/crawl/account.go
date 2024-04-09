package crawl

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/stonecool/1701livehouse-server/internal/model"
	"github.com/stonecool/1701livehouse-server/pkg/cache"
	"log"
	"reflect"
)

type Account struct {
	ID          int                    `json:"id"`
	TemplateId  string                 `json:"template_id"`
	Name        string                 `json:"name"`
	Username    string                 `json:"username"`
	Password    string                 `json:"password"`
	Headers     map[string]interface{} `json:"headers"`
	QueryParams string                 `json:"query_params"`
	FormData    string                 `json:"form_data"`
	State       uint8                  `json:"state"`

	Collector *colly.Collector
	Ch        chan CmdRequest
}

var accountCache *cache.Memo

func init() {
	accountCache = cache.New(getCrawlAccountByID)
}

func (a *Account) Add() error {
	headers, err := json.Marshal(a.Headers)
	if err != nil {
		log.Printf("json marshal error: %s", err)
		return err
	}

	account := map[string]interface{}{
		"template_id":  a.TemplateId,
		"name":         a.Name,
		"username":     a.Username,
		"password":     a.Password,
		"headers":      string(headers),
		"query_params": a.QueryParams,
		"form_data":    a.FormData,
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
	if err := json.Unmarshal([]byte(a.QueryParams), &params); err != nil {
		return nil, err
	}

	return params, nil
}

// GetFormData
func (a *Account) GetFormData() string {

	return a.FormData
}

func GetCrawlAccountByID(id int) (*Account, error) {
	if t, err := accountCache.Get(id); err != nil {
		return nil, err
	} else {
		return t.(*Account), nil
	}
}

func getCrawlAccountByID(id int) (interface{}, error) {
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
		ID:        p.ID,
		Name:      p.Name,
		Headers:   make(map[string]interface{}),
		Collector: colly.NewCollector(),
		Ch:        make(chan CmdRequest),
	}

	return &account, nil
}

func (a *Account) GetId() int {
	return a.ID
}

func (a *Account) GetChan() chan CmdRequest {
	return a.Ch
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
