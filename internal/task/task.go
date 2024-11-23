package task

type Task struct {
	ID        int    `json:"id"`
	Category  string `json:"category"`
	TargetId  string `json:"target_id"`
	MetaType  string `json:"meta_type"`
	MetaId    int    `json:"meta_id"`
	CronSpec  string
	FirstTime int `json:"first_time"`
	LastTime  int `json:"last_time"`
	Count     int `json:"count"`
	mark      string
}

func (t *Task) GetId() int {
	return t.ID
}
