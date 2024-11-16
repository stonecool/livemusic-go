package database

import (
	"gorm.io/gorm"
)

type Musician struct {
	BaseModel

	Name string `json:"name"`
}

// AddMusician Adds a musician
func AddMusician(data map[string]interface{}) error {
	musician := Musician{
		Name: data["name"].(string),
	}

	if err := DB.Create(&musician).Error; err != nil {
		return err
	}

	return nil
}

// DeleteMusician Deletes a musician based on ID
func DeleteMusician(id int) error {
	if err := DB.Where("id = ?", id).Delete(Musician{}).Error; err != nil {
		return err
	}

	return nil
}

// EditMusician Edits a musician based on ID
func EditMusician(id int, data interface{}) error {
	if err := DB.Model(&Musician{}).Where("id = ? AND deleted_at = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// GetMusician Gets a musician based on ID
func GetMusician(id int) (*Musician, error) {
	var musician Musician
	if err := DB.Where("id = ? AND deleted_at = ?", id, 0).First(&musician).Error; err != nil && err != gorm.ErrNotImplemented {
		return nil, err
	}

	// TODO related table
	//err = DB.Model(&article).Related(&article.Tag).Error
	//if err != nil && err != gorm.ErrRecordNotFound {
	//	return nil, err
	//}

	return &musician, nil
}

// ExistMusicianById checks if a musician exists based on ID
func ExistMusicianById(id int) (bool, error) {
	var musician Musician
	if err := DB.Select("id").Where("id = ?", id).First(&musician).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return musician.ID > 0, nil
}

// GetMusicianCount Gets count of musician based on maps condition
func GetMusicianCount(maps interface{}) (int64, error) {
	var count int64
	// TODO Preload
	if err := DB.Model(&Musician{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetMusicians Gets musicians on page
func GetMusicians(pageNum int, pageSize int, maps interface{}) ([]*Musician, error) {
	var musicians []*Musician
	// TODO Preload
	//if err := DB.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&musicians).Error; err != nil {
	//	return nil, err
	//}

	return musicians, nil
}
