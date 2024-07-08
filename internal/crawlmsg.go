package internal

import (
	"github.com/stonecool/livemusic-go/internal/model"
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

// AddCrawlMsg
func AddCrawlMsg(dataType string, dataId int, crawlType string, accountId string) (*CrawlMsg, error) {
	_, ok := CrawlAccountMap[crawlType]
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
		msg := CrawlMsg{
			ID:        m.ID,
			DataType:  m.DataType,
			DataId:    m.DataId,
			CrawlType: m.CrawlType,
			AccountId: m.AccountId,
			Count:     m.Count,
			FirstTime: m.FirstTime,
			LastTime:  m.LastTime,
			mark:      m.Mark,
		}

		return &msg, nil
	}
}
