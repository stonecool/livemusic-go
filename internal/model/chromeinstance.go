package model

import (
	"gorm.io/gorm"
)

type ChromeInstance struct {
	Model

	IP          string
	Port        int
	DebuggerUrl string
	Status      int
}

func AddChromeInstance(data map[string]interface{}) (*ChromeInstance, error) {
	ins := ChromeInstance{
		IP:          data["ip"].(string),
		Port:        data["port"].(int),
		DebuggerUrl: data["debugger_url"].(string),
		Status:      data["status"].(int),
	}

	if err := DB.Create(&ins).Error; err != nil {
		return nil, err
	}

	return &ins, nil
}

func ExistsChromeInstance(ip string, port int) (bool, error) {
	var ins ChromeInstance
	err := DB.Select("id").Where("ip = '?' AND port = ?",
		ip, port).First(&ins).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return ins.ID > 0, nil
}

func GetChromeInstance(id int) (*ChromeInstance, error) {
	var instance ChromeInstance
	if err := DB.Where("id = ? AND deleted_at = ?", id, 0).First(&instance).Error; err != nil {
		return nil, err
	} else {
		return &instance, err
	}
}

func GetChromeInstanceAll() ([]*ChromeInstance, error) {
	var instances []*ChromeInstance
	if err := DB.Where("deleted_at = ?", 0).Find(&instances).Error; err != nil {
		return nil, err
	}

	return instances, nil
}
