package internal

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
	"reflect"
)

type CrawlRoutine struct {
	ID              int    `json:"id"`
	DataType        string `json:"data_type"`
	DataId          int    `json:"data_id"`
	AccountType     string `json:"account_type"`
	TargetAccountId string `json:"target_account_id"`
	Count           int    `json:"count"`
	FirstTime       int    `json:"first_time"`
	LastTime        int    `json:"last_time"`
	mark            string
}

func (r *CrawlRoutine) init(msg *model.CrawlRoutine) {
	r.ID = msg.ID
	r.DataType = msg.DataType
	r.DataId = msg.DataId
	r.AccountType = msg.AccountType
	r.TargetAccountId = msg.TargetAccountId
	r.Count = msg.Count
	r.FirstTime = msg.FirstTime
	r.LastTime = msg.LastTime
	r.mark = msg.Mark
}

func (r *CrawlRoutine) Add() error {
	_, ok := config.AccountMap[r.AccountType]
	if !ok {
		return fmt.Errorf("account_type:%s not exists", r.AccountType)
	}

	exist, err := dataTypeIdExists(r.DataType, r.DataId)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("data table not exists")
	}

	if exist, err := model.ExistCrawlRoutine(r.DataType, r.DataId, r.AccountType); err != nil {
		Logger.Warn("routine exists")

		return fmt.Errorf("some error")
	} else {
		if exist {
			return fmt.Errorf("exists")
		}
	}

	data := map[string]interface{}{
		"data_type":         r.DataType,
		"data_id":           r.DataId,
		"account_type":      r.AccountType,
		"target_account_id": r.TargetAccountId,
	}

	if routine, err := model.AddCrawlRoutine(data); err != nil {
		return err
	} else {
		r.init(routine)
		return nil
	}
}

func (r *CrawlRoutine) Get() error {
	if routine, err := model.GetCrawlRoutine(r.ID); err != nil {
		return err
	} else {
		r.init(routine)
		return nil
	}
}

func (r *CrawlRoutine) GetAll() ([]*CrawlRoutine, error) {
	if s, err := model.GetCrawlRoutineAll(); err != nil {
		return nil, err
	} else {
		var ret []*CrawlRoutine

		for _, routine := range s {
			tempMsg := &CrawlRoutine{}
			tempMsg.init(routine)
			ret = append(ret, tempMsg)
		}

		return ret, nil
	}
}

func (r *CrawlRoutine) Delete() error {
	routine, err := model.GetCrawlRoutine(r.ID)
	if err != nil {
		return err
	}

	return model.DeleteCrawlRoutine(routine)
}

func (r *CrawlRoutine) Edit() error {
	routine, err := model.GetCrawlRoutine(r.ID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"data_type":         r.DataType,
		"data_id":           r.DataId,
		"account_type":      r.AccountType,
		"target_account_id": r.TargetAccountId,
	}

	routine, err = model.EditCrawlRoutine(r.ID, data)
	if err != nil {
		return err
	} else {
		r.init(routine)
		return nil
	}
}

func dataTypeIdExists(dataType string, dataId int) (bool, error) {
	val, ok := dataType2StructMap[dataType]
	if !ok {
		return false, fmt.Errorf("data_type:%s illegal", dataType)
	}

	originalType := reflect.TypeOf(val).Elem()
	newVar := reflect.New(originalType).Elem()

	pointer := newVar.Addr().Interface().(IDataTable)
	pointer.setId(dataId)

	return pointer.exist()
}
