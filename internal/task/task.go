package task

type Task struct {
	ID        int    `json:"id"`
	Category  string `json:"category"`
	TargetID  string `json:"target_id"`
	MetaType  string `json:"meta_type"`
	MetaID    int    `json:"meta_id"`
	Count     int    `json:"count"`
	FirstTime int    `json:"first_time"`
	LastTime  int    `json:"last_time"`
	mark     string
	CronSpec string
}

func NewTask(m *Task) *Task {
	return &Task{
		ID:        m.ID,
		Category:  m.Category,
		TargetID:  m.TargetID,
		MetaType:  m.MetaType,
		MetaID:    m.MetaID,
		Count:     m.Count,
		FirstTime: m.FirstTime,
		LastTime:  m.LastTime,
		mark:      m.mark,
		CronSpec:  m.CronSpec,
	}
}


