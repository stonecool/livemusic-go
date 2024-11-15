package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/stonecool/livemusic-go/internal/chrome"
	"github.com/stonecool/livemusic-go/internal/config"
	"github.com/stonecool/livemusic-go/internal/model"
)

type CrawlTask struct {
	ID              int    `json:"id"`
	DataType        string `json:"data_type"`
	DataId          int    `json:"data_id"`
	AccountType     string `json:"account_type"`
	TargetAccountId string `json:"target_account_id"`
	Count           int    `json:"count"`
	FirstTime       int    `json:"first_time"`
	LastTime        int    `json:"last_time"`
	CronSpec        string `json:"cron_spec"`
	mark            string
}

func (t *CrawlTask) init(task *model.CrawlTask) {
	t.ID = task.ID
	t.DataType = task.DataType
	t.DataId = task.DataId
	t.AccountType = task.AccountType
	t.TargetAccountId = task.TargetAccountId
	t.Count = task.Count
	t.FirstTime = task.FirstTime
	t.LastTime = task.LastTime
	t.CronSpec = task.CronSpec
	t.mark = task.Mark
}

func (t *CrawlTask) Add() error {
	_, ok := config.AccountMap[t.AccountType]
	if !ok {
		return fmt.Errorf("account_type:%s not exists", t.AccountType)
	}

	exist, err := dataTypeIdExists(t.DataType, t.DataId)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("data table not exists")
	}

	if exist, err := model.ExistCrawlTask(t.DataType, t.DataId, t.AccountType); err != nil {
		Logger.Warn("task exists")
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

	if task, err := model.AddCrawlTask(data); err != nil {
		return err
	} else {
		t.init(task)
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
