package internal

import (
	"fmt"
	"log"
	reflect "reflect"
	"time"

	"github.com/stonecool/livemusic-go/internal/chrome"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
)

type CrawlTask struct {
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

func NewCrawlTask(m *model.CrawlTask) *CrawlTask {
	return &CrawlTask{
		ID:        m.ID,
		Category:  m.Category,
		TargetID:  m.TargetId,
		MetaType:  m.MetaType,
		MetaID:    m.MetaID,
		Count:     m.Count,
		FirstTime: m.FirstTime,
		LastTime:  m.LastTime,
		mark:      m.Mark,
		cronSpec:  m.CronSpec,
	}
}

func (t *CrawlTask) Add() error {
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

	if exist, err := model.ExistCrawlTask(t.MetaType, t.MetaID, t.Category); err != nil {
		Logger.Warn("m exists")
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

	if m, err := model.AddCrawlTask(data); err != nil {
		return err
	} else {
		t := NewCrawlTask(m)
		// 添加到调度器
		return GetScheduler().AddTask(t)
	}
}

func (t *CrawlTask) Execute() error {
	const (
		maxRetries = 3 // 最大重试次数
		retryDelay = 5 // 重试间隔(秒)
	)

	msg := NewAsyncMessage(&Message{
		Cmd:  CrawlCmd_Crawl,
		Data: t,
	})

	task := &CrawlTask{
		Category: t.AccountType,
		Message:  msg,
	}

	var lastErr error
	for retry := 0; retry < maxRetries; retry++ {
		// 如果不是第一次尝试,等待一段时间
		if retry > 0 {
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}

		// 尝试分发任务
		if err := chrome.GetPool().DispatchTask(t.AccountType, task); err == nil {
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

func GetAllCrawlTasks() ([]*CrawlTask, error) {
	modelTasks, err := model.GetCrawlTaskAll()
	if err != nil {
		return nil, err
	}

	var tasks []*CrawlTask
	for _, mt := range modelTasks {
		task := &CrawlTask{}
		task.init(mt)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func dataTypeIdExists(dataType string, dataId int) (bool, error) {
	val, ok := dataType2StructMap[dataType]
	if !ok {
		return false, fmt.Errorf("data_type:%s illegal", dataType)
	}

	originalType := reflect.TypeOf(val).Elem()
	newVar := reflect.New(originalType).Elem()

	pointer := newVar.Addr().Interface().(IDataTable)
	pointer.setId(dataId)

	return pointer.exist()
}
