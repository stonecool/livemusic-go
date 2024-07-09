package internal

import (
	"github.com/stonecool/livemusic-go/internal/model"
	"reflect"
)

type MsgProducer struct {
	ID        int    `json:"id"`
	DataType  string `json:"data_type"`
	DataId    int    `json:"data_id"`
	CrawlType string `json:"crawl_type"`
	AccountId string `json:"account_id"`
	Count     int    `json:"count"`
	FirstTime int    `json:"first_time"`
	LastTime  int    `json:"last_time"`
	mark      string
}

func (msgProducer *MsgProducer) init(m *model.MsgProducer) {
	if m == nil || reflect.ValueOf(m).IsZero() {
		return
	}

	msgProducer.ID = m.ID
	msgProducer.DataType = m.DataType
	msgProducer.DataId = m.DataId
	msgProducer.CrawlType = m.CrawlType
	msgProducer.AccountId = m.AccountId
	msgProducer.Count = m.Count
	msgProducer.FirstTime = m.FirstTime
	msgProducer.LastTime = m.LastTime
	msgProducer.mark = m.Mark
}

// AddMsgProducer
func AddMsgProducer(dataType string, dataId int, crawlType string, accountId string) (*MsgProducer, error) {
	_, ok := CrawlAccountMap[crawlType]
	if !ok {
		return nil, error(nil)
	}

	if model.MsgProducerExists(dataType, dataId, crawlType) {
		Logger.Warn("coroutine exists")

		return nil, error(nil)
	}

	data := map[string]interface{}{
		"data_type":  dataType,
		"data_id":    dataId,
		"crawl_type": crawlType,
		"account_id": accountId,
	}

	if m, err := model.AddMsgProducer(data); err != nil {
		return nil, err
	} else {
		msgProducer := MsgProducer{}
		msgProducer.init(m)

		return &msgProducer, nil
	}
}
