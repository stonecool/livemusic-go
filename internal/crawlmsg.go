package internal

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
	"reflect"
)

type CrawlMsg struct {
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

func (m *CrawlMsg) init(msg *model.CrawlMsg) {
	m.ID = msg.ID
	m.DataType = msg.DataType
	m.DataId = msg.DataId
	m.AccountType = msg.AccountType
	m.TargetAccountId = msg.TargetAccountId
	m.Count = msg.Count
	m.FirstTime = msg.FirstTime
	m.LastTime = msg.LastTime
	m.mark = msg.Mark
}

func (m *CrawlMsg) Add() error {
	_, ok := config.AccountMap[m.AccountType]
	if !ok {
		return fmt.Errorf("account_type:%s not exists", m.AccountType)
	}

	exist, err := dataTypeIdExists(m.DataType, m.DataId)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("data table not exists")
	}

	if exist, err := model.ExistCrawlMsg(m.DataType, m.DataId, m.AccountType); err != nil {
		Logger.Warn("msg exists")

		return fmt.Errorf("some error")
	} else {
		if exist {
			return fmt.Errorf("exists")
		}
	}

	data := map[string]interface{}{
		"data_type":         m.DataType,
		"data_id":           m.DataId,
		"account_type":      m.AccountType,
		"target_account_id": m.TargetAccountId,
	}

	if msg, err := model.AddCrawlMsg(data); err != nil {
		return err
	} else {
		m.init(msg)
		return nil
	}
}

func (m *CrawlMsg) Get() error {
	if msg, err := model.GetCrawlMg(m.ID); err != nil {
		return err
	} else {
		m.init(msg)
		return nil
	}
}

func (m *CrawlMsg) GetAll() ([]*CrawlMsg, error) {
	if msgs, err := model.GetCrawlMsgAll(); err != nil {
		return nil, err
	} else {
		var s []*CrawlMsg

		for _, msg := range msgs {
			tempMsg := &CrawlMsg{}
			tempMsg.init(msg)
			s = append(s, tempMsg)
		}

		return s, nil
	}
}

func (m *CrawlMsg) Delete() error {
	msg, err := model.GetCrawlMg(m.ID)
	if err != nil {
		return err
	}

	return model.DeleteCrawlMsg(msg)
}

func (m *CrawlMsg) Edit() error {
	msg, err := model.GetCrawlMg(m.ID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"data_type":         m.DataType,
		"data_id":           m.DataId,
		"account_type":      m.AccountType,
		"target_account_id": m.TargetAccountId,
	}

	msg, err = model.EditCrawlMsg(m.ID, data)
	if err != nil {
		return err
	} else {
		m.init(msg)
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
