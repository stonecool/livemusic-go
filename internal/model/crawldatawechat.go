package model

type CrawlDataWechat struct {
	Model

	RID          int `gorm:"column:rid"`
	UID          string
	Title        string
	Cover        string
	Link         string
	OriginalLink string
	Datetime     int
	State        uint8
	DataType     uint8 `gorm:"column:type"`
	DataSubType  uint8 `gorm:"column:sub_type"`
	RawId        int
}

type CrawlDataWechatRaw struct {
	CrawlDataWechat
}

func addCrawlDataWechat(data map[string]interface{}, tableName string) error {
	m := CrawlDataWechat{
		RID:          data["rid"].(int),
		UID:          data["uid"].(string),
		Title:        data["title"].(string),
		Cover:        data["cover"].(string),
		Link:         data["link"].(string),
		OriginalLink: data["original_link"].(string),
		Datetime:     data["datetime"].(int),
		State:        data["state"].(uint8),
		DataType:     data["data_type"].(uint8),
		DataSubType:  data["data_sub_type"].(uint8),
		RawId:        data["raw_id"].(int),
	}

	return DB.Table(tableName).Create(&m).Error
}

func AddCrawlDataWechat(data map[string]interface{}) error {
	return addCrawlDataWechat(data, "crawl_data_wechat")
}

func AddCrawlDataWechatRaw(data map[string]interface{}) error {
	return addCrawlDataWechat(data, "crawl_data_wechat_raw")
}
