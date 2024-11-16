package chrome

import (
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type model struct {
	database.BaseModel

	IP          string
	Port        int
	DebuggerURL string
	State       int
}

func AddChromeInstance(data map[string]interface{}) (*model, error) {
	ins := model{
		IP:          data["ip"].(string),
		Port:        data["port"].(int),
		DebuggerURL: data["debugger_url"].(string),
		State:       data["state"].(int),
	}

	if err := DB.Create(&ins).Error; err != nil {
		return nil, err
	}

	return &ins, nil
}

func ExistsChromeInstance(ip string, port int) (bool, error) {
	var ins model
	err := DB.Select("id").Where("ip = '?' AND port = ?",
		ip, port).First(&ins).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return ins.ID > 0, nil
}

func GetChromeInstance(id int) (*model, error) {
	var instance model
	if err := DB.Where("id = ? AND deleted_at = ?", id, 0).First(&instance).Error; err != nil {
		return nil, err
	} else {
		return &instance, err
	}
}

func GetChromeInstanceAll() ([]*model, error) {
	var instances []*model
	if err := DB.Where("deleted_at = ?", 0).Find(&instances).Error; err != nil {
		return nil, err
	}

	return instances, nil
}
