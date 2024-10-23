package model

type ChromeInstance struct {
	Model

	Addr  string
	State int
}

func AddChromeInstance(data map[string]interface{}) (*ChromeInstance, error) {
	ins := ChromeInstance{
		Addr:  data["addr"].(string),
		State: data["state"].(int),
	}

	if err := db.Create(&ins).Error; err != nil {
		return nil, err
	}

	return &ins, nil
}
