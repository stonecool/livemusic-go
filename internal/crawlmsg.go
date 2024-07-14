package internal

import (
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
	"reflect"
)

type CrawlMsg struct {
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

func initCrawlMsg(m *model.CrawlMsg) *CrawlMsg {
	if m == nil || reflect.ValueOf(m).IsZero() {
		return nil
	}

	var msg CrawlMsg
	msg.ID = m.ID
	msg.DataType = m.DataType
	msg.DataId = m.DataId
	msg.CrawlType = m.CrawlType
	msg.AccountId = m.AccountId
	msg.Count = m.Count
	msg.FirstTime = m.FirstTime
	msg.LastTime = m.LastTime
	msg.mark = m.Mark

	return &msg
}

// AddCrawlMsg
func AddCrawlMsg(dataType string, dataId int, crawlType string, accountId string) (*CrawlMsg, error) {
	_, ok := config.AccountMap[crawlType]
	if !ok {
		return nil, error(nil)
	}

	if model.CrawlMsgExists(dataType, dataId, crawlType) {
		Logger.Warn("coroutine exists")

		return nil, error(nil)
	}

	data := map[string]interface{}{
		"data_type":  dataType,
		"data_id":    dataId,
		"crawl_type": crawlType,
		"account_id": accountId,
	}

	if m, err := model.AddCrawlMsg(data); err != nil {
		return nil, err
	} else {
		return initCrawlMsg(m), nil
	}
}

func GetCrawlMsg(id int) (*CrawlMsg, error) {
	msg, err := model.GetCrawlMg(id)
	if err != nil {
		return nil, err
	} else {
		return initCrawlMsg(msg), nil
	}
}
