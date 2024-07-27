package model

import (
	"gorm.io/gorm"
	"time"
)

type Livehouse struct {
	Model

	Name string `json:"name"`
}

// AddLiveHouse Adds a live house
func AddLiveHouse(data map[string]interface{}) (*Livehouse, error) {
	house := Livehouse{
		Name: data["name"].(string),
	}

	if err := db.Create(&house).Error; err != nil {
		return nil, err
	}

	return &house, nil
}

// GetLiveHouse gets a live house based on ID
func GetLiveHouse(id int) (*Livehouse, error) {
	var liveHouse Livehouse
	if err := db.Where("id = ? AND deleted_at = ?", id, 0).First(&liveHouse).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &liveHouse, nil
}

// GetLiveHouses Gets all livehouses
func GetLiveHouses() ([]*Livehouse, error) {
	var liveHouses []*Livehouse
	if err := db.Where("deleted_at = ?", 0).Find(&liveHouses).Error; err != nil {
		return nil, err
	}

	return liveHouses, nil
}

// EditLiveHouse edits a live house based on ID
func EditLiveHouse(id int, data interface{}) error {
	if err := db.Model(&Livehouse{}).Where("id = ? AND deleted_at = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// DeleteLiveHouse deletes a live house based on id
func DeleteLiveHouse(id int) error {
	return db.Model(&Livehouse{}).Where("id = ? AND deleted_at = ?", id, 0).Update("deleted_at", time.Now().Unix()).Error
}

func ExistLivehouse(id int) (bool, error) {
	var livehouse Livehouse
	err := db.Select("id").Where("id = ? AND deleted_at = ?", id, 0).First(&livehouse).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return livehouse.ID > 0, nil
}
