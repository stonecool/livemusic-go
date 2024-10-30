package model

import "gorm.io/gorm"

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

	if err := db.Create(&ins).Error; err != nil {
		return nil, err
	}

	return &ins, nil
}

func ExistsChromeInstance(ip string, port int) (bool, error) {
	var ins ChromeInstance
	err := db.Select("id").Where("ip = '?' AND port = ?",
		ip, port).First(&ins).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return ins.ID > 0, nil
}