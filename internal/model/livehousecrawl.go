package model

import (
	"gorm.io/gorm"
)

type LiveHouseCrawl struct {
	Model

	Name       string
	Type       int8
	URL        string
	Headers    []byte
	CrawlCount int32
	FirstTime  int32
	LastTime   int32
	Status     int8
}

// AddLiveHouseCrawlMeta a live house server
func AddLiveHouseCrawlMeta(data map[string]interface{}) error {
	// TODO (string)
	crawlMeta := LiveHouseCrawl{
		Name:       data["name"].(string),
		Type:       int8(data["type"].(int)),
		URL:        data["url"].(string),
		Headers:    data["headers"].([]byte),
		CrawlCount: int32(data["crawlCount"].(int)),
		FirstTime:  int32(data["firstTime"].(int)),
		LastTime:   int32(data["lastTime"].(int)),
		Status:     int8(data["status"].(int)),
	}

	if err := db.Create(&crawlMeta).Error; err != nil {
		return err
	}

	return nil
}

// GetLiveHouseCrawlMeta Get a single article based on id
func GetLiveHouseCrawlMeta(id int) (*LiveHouseCrawl, error) {
	var crawlMeta LiveHouseCrawl
	if err := db.Where("id = ?", id).First(&crawlMeta).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &crawlMeta, nil
}
