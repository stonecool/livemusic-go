package database

import (
	"gorm.io/gorm"
	"time"
)

type Livehouse struct {
	BaseModel

	Name string `json:"name"`
}

// AddLiveHouse Adds a live house
func AddLiveHouse(data map[string]interface{}) (*Livehouse, error) {
	house := Livehouse{
		Name: data["name"].(string),
	}

	if err := DB.Create(&house).Error; err != nil {
		return nil, err
	}

	return &house, nil
}

// GetLiveHouse gets a live house based on ID
func GetLiveHouse(id int) (*Livehouse, error) {
	var liveHouse Livehouse
	if err := DB.Where("id = ? AND deleted_at = ?", id, 0).First(&liveHouse).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &liveHouse, nil
}

// GetLiveHouses Gets all livehouses
func GetLiveHouses() ([]*Livehouse, error) {
	var liveHouses []*Livehouse
	if err := DB.Where("deleted_at = ?", 0).Find(&liveHouses).Error; err != nil {
		return nil, err
	}

	return liveHouses, nil
}

// EditLiveHouse edits a live house based on ID
func EditLiveHouse(id int, data interface{}) error {
	return DB.Model(&Livehouse{}).Where("id = ? AND deleted_at = ?", id, 0).Updates(data).Error
}

// DeleteLiveHouse deletes a live house based on id
func DeleteLiveHouse(id int) error {
	return DB.Model(&Livehouse{}).Where("id = ? AND deleted_at = ?", id, 0).Update("deleted_at", time.Now().Unix()).Error
}

func ExistLivehouse(id int) (bool, error) {
	var livehouse Livehouse
	err := DB.Select("id").Where("id = ? AND deleted_at = ?", id, 0).First(&livehouse).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return livehouse.ID > 0, nil
}
