package model

import (
	"gorm.io/gorm"
)

type LiveHouse struct {
	Model

	Name string `json:"name"`
}

// AddLiveHouse Adds a live house
func AddLiveHouse(data map[string]interface{}) error {
	liveHouse := LiveHouse{
		Name: data["name"].(string),
	}

	if err := db.Create(&liveHouse).Error; err != nil {
		return err
	}

	return nil
}

// DeleteLiveHouse deletes a live house based on id
func DeleteLiveHouse(id int) error {
	if err := db.Where("id = ?", id).Delete(LiveHouse{}).Error; err != nil {
		return err
	}

	return nil
}

// EditLiveHouse edits a live house based on ID
func EditLiveHouse(id int, data interface{}) error {
	if err := db.Model(&LiveHouse{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// GetLiveHouse gets a live house based on ID
func GetLiveHouse(id int) (*LiveHouse, error) {
	var liveHouse LiveHouse
	if err := db.Where("id = ? AND deleted_on = ?", id, 0).First(&liveHouse).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &liveHouse, nil
}

// ExistsMusicianById checks if a live house exists based on ID
func ExistsMusicianById(id int) (bool, error) {
	var liveHouse LiveHouse
	if err := db.Select("id").Where("id = ?", id).First(&liveHouse).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return liveHouse.ID > 0, nil
}

// GetLiveHouseCount Gets count of live house based on maps condition
func GetLiveHouseCount(maps interface{}) (int64, error) {
	var count int64

	if err := db.Model(&LiveHouse{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetLiveHouses Gets live houses on page
func GetLiveHouses(pageNum int, pageSize int, maps interface{}) ([]*LiveHouse, error) {
	var liveHouses []*LiveHouse

	return liveHouses, nil
}
