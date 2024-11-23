package task

type Task struct {
	ID        int    `json:"id"`
	Category  string `json:"category"`
	TargetID  string `json:"target_id"`
	MetaType  string `json:"meta_type"`
	MetaID    int    `json:"meta_id"`
	CronSpec  string
	FirstTime int `json:"first_time"`
	LastTime  int `json:"last_time"`
	Count     int `json:"count"`
	mark      string
}

func (t *Task) GetID() int {
	return t.ID
}

func (t *Task) GetCategory() string {
	return t.Category
}

func (t *Task) Execute() error {
	return nil
}

func (t *Task) Cancel() error {
	return nil
}

func (t *Task) GetCronSpec() string {
	return t.CronSpec
}
