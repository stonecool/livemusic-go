package instance

import (
	"github.com/stonecool/livemusic-go/internal/model"
	"gorm.io/gorm"
)

type instanceModel struct {
	model.Model

	IP          string
	Port        int
	DebuggerUrl string
	Status      int
}

func AddChromeInstance(data map[string]interface{}) (*instanceModel, error) {
	ins := instanceModel{
		IP:          data["ip"].(string),
		Port:        data["port"].(int),
		DebuggerUrl: data["debugger_url"].(string),
		Status:      data["status"].(int),
	}

	if err := model.DB.Create(&ins).Error; err != nil {
		return nil, err
	}

	return &ins, nil
}

func ExistsChromeInstance(ip string, port int) (bool, error) {
	var ins instanceModel
	err := model.DB.Select("id").Where("ip = '?' AND port = ?",
		ip, port).First(&ins).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return ins.ID > 0, nil
}
