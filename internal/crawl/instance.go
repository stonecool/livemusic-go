package crawl

import (
	"encoding/json"
	"github.com/stonecool/1701livehouse-server/internal/model"
	"github.com/stonecool/1701livehouse-server/pkg/cache"
	"log"
	"reflect"
)

type Instance struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	AccountId   int                    `json:"account_id"`
	Headers     map[string]interface{} `json:"headers"`
	QueryParams string                 `json:"query_params"`
	FormData    string                 `json:"form_data"`
	State       uint                   `json:"state"`
	Count       int                    `json:"count"`
	FirstTime   int                    `json:"first_time"`
	LastTime    int                    `json:"last_time"`
}

var instanceCache *cache.Memo

func init() {
	instanceCache = cache.New(getCrawlInstanceByID)
}

func (i *Instance) Add() error {
	instance := map[string]interface{}{
		"name":       i.Name,
		"account_id": i.AccountId,
		"state":      i.State,
		"count":      i.Count,
	}

	if id, err := model.AddCrawlInstance(instance); err != nil {
		return err
	} else {
		i.ID = id
		return nil
	}
}

// GetHeaders
func (i *Instance) GetHeaders() (map[string]interface{}, error) {
	return i.Headers, nil
}

// GetQueryParams
func (i *Instance) GetQueryParams() (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(i.QueryParams), &params); err != nil {
		return nil, err
	}

	return params, nil
}

// GetFormData
func (i *Instance) GetFormData() string {
	return i.FormData
}

func GetCrawlInstanceByID(id int) (*Instance, error) {
	if i, err := instanceCache.Get(id); err != nil {
		return nil, err
	} else {
		return i.(*Instance), nil
	}
}

func getCrawlInstanceByID(id int) (interface{}, error) {
	var instance *model.CrawlInstance
	instance, err := model.GetCrawlInstance(id)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	if reflect.ValueOf(*instance).IsZero() {
		return &Account{}, nil
	}

	newInstance := Instance{
		ID:        instance.ID,
		Name:      instance.Name,
		AccountId: instance.AccountId,
	}

	return &newInstance, nil
}
