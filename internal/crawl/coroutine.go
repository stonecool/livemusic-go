package crawl

import (
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/model"
)

type Coroutine struct {
	ID        int    `json:"id"`
	DataType  string `json:"data_type"`
	DataId    int    `json:"data_id"`
	CrawlType string `json:"crawl_type"`
	AccountId string `json:"account_id"`
	Count     int    `json:"count"`
	FirstTime int    `json:"first_time"`
	LastTime  int    `json:"last_time"`
	crawlMark string
}

// AddCoroutine
func AddCoroutine(dataType string, dataId int, crawlType string, accountId string) (*Coroutine, error) {
	_, ok := internal.CrawlAccountMap[crawlType]
	if !ok {
		return nil, error(nil)
	}

	if model.CrawlCoroutineExists(dataType, dataId, crawlType) {
		internal.Logger.Warn("coroutine exists")

		return nil, error(nil)
	}

	data := map[string]interface{}{
		"data_type":  dataType,
		"data_id":    dataId,
		"crawl_type": crawlType,
		"account_id": accountId,
	}

	if m, err := model.AddCrawlCoroutine(data); err != nil {
		return nil, err
	} else {
		coroutine := Coroutine{
			ID:        m.ID,
			CrawlType: m.CrawlType,
		}

		return &coroutine, nil
	}
}
