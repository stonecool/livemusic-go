package internal

import (
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
)

type CrawlMsg struct {
	ID          int    `json:"id"`
	DataType    string `json:"data_type"`
	DataId      int    `json:"data_id"`
	AccountType string `json:"account_type"`
	AccountId   string `json:"account_id"`
	Count       int    `json:"count"`
	FirstTime   int    `json:"first_time"`
	LastTime    int    `json:"last_time"`
	mark        string
}

func (m *CrawlMsg) init(msg *model.CrawlMsg) {
	m.ID = msg.ID
	m.DataType = msg.DataType
	m.DataId = msg.DataId
	m.AccountType = msg.AccountType
	m.AccountId = msg.AccountId
	m.Count = msg.Count
	m.FirstTime = msg.FirstTime
	m.LastTime = msg.LastTime
	m.mark = msg.Mark
}

func (m *CrawlMsg) Add() error {
	_, ok := config.AccountMap[m.AccountType]
	if !ok {
		return error(nil)
	}

	if model.CrawlMsgExists(m.DataType, m.DataId, m.AccountType) {
		Logger.Warn("msg exists")

		return error(nil)
	}

	data := map[string]interface{}{
		"data_type":    m.DataType,
		"data_id":      m.DataId,
		"account_type": m.AccountType,
		"account_id":   m.AccountId,
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
		"data_type":    m.DataType,
		"data_id":      m.DataId,
		"account_type": m.AccountType,
		"account_id":   m.AccountId,
	}

	msg, err = model.EditCrawlMsg(m.ID, data)
	if err != nil {
		return err
	} else {
		m.init(msg)
		return nil
	}
}
