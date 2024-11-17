package task

import (
	"fmt"
	"log"
	"time"

	"github.com/stonecool/livemusic-go/internal/chrome"
)

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

func (t *Task) Execute() error {
	const (
		maxRetries = 3 // 最大重试次数
		retryDelay = 5 // 重试间隔(秒)
	)

	//msg := client.NewAsyncMessage(&internal.Message{
	//	Cmd:  internal.CrawlCmd_Crawl,
	//	Data: t,
	//})

	task := &Task{
		Category: t.Category,
		//Message:  msg,
	}

	var lastErr error
	for retry := 0; retry < maxRetries; retry++ {
		// 如果不是第一次尝试,等待一段时间
		if retry > 0 {
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}

		// 尝试分发任务
		if err := chrome.GetPool().DispatchTask(t.Category, task); err == nil {
			return nil
		} else {
			lastErr = err
			log.Printf("Task %d dispatch failed (attempt %d/%d): %v",
				t.ID, retry+1, maxRetries, err)
		}
	}

	return fmt.Errorf("dispatch task failed after %d attempts: %v",
		maxRetries, lastErr)
}


