package task

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/client"
	"github.com/stonecool/livemusic-go/internal/scheduler"
	"log"
	reflect "reflect"
	"time"

	"github.com/stonecool/livemusic-go/internal/chrome"
	"github.com/stonecool/livemusic-go/internal/config"
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
	mark      string
	cronSpec  string
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
		cronSpec:  m.cronSpec,
	}
}

func (t *Task) Add() error {
	_, ok := config.AccountMap[t.Category]
	if !ok {
		return fmt.Errorf("account_type:%s not exists", t.Category)
	}

	exist, err := dataTypeIdExists(t.MetaType, t.MetaID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("data table not exists")
	}

	if exist, err := ExistCrawlTask(t.MetaType, t.MetaID, t.Category); err != nil {
		internal.Logger.Warn("m exists")
		return fmt.Errorf("some error")
	} else if exist {
		return fmt.Errorf("exists")
	}

	data := map[string]interface{}{
		"data_type":         t.DataType,
		"data_id":           t.DataId,
		"account_type":      t.AccountType,
		"target_account_id": t.TargetAccountId,
		"cron_spec":         t.CronSpec,
	}

	if m, err := AddCrawlTask(data); err != nil {
		return err
	} else {
		t := NewTask(m)
		// 添加到调度器
		return scheduler.GetScheduler().AddTask(t)
	}
}

func (t *Task) Execute() error {
	const (
		maxRetries = 3 // 最大重试次数
		retryDelay = 5 // 重试间隔(秒)
	)

	msg := client.NewAsyncMessage(&internal.Message{
		Cmd:  internal.CrawlCmd_Crawl,
		Data: t,
	})

	task := &Task{
		Category: t.Category,
		Message:  msg,
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

func GetAllCrawlTasks() ([]*Task, error) {
	modelTasks, err := GetCrawlTaskAll()
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for _, mt := range modelTasks {
		task := &Task{}
		task.init(mt)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func dataTypeIdExists(dataType string, dataId int) (bool, error) {
	val, ok := internal.DataType2StructMap[dataType]
	if !ok {
		return false, fmt.Errorf("data_type:%s illegal", dataType)
	}

	originalType := reflect.TypeOf(val).Elem()
	newVar := reflect.New(originalType).Elem()

	pointer := newVar.Addr().Interface().(internal.IDataTable)
	pointer.setId(dataId)

	return pointer.exist()
}
