package scheduler

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/stonecool/livemusic-go/internal/chrome/pool"
	"github.com/stonecool/livemusic-go/internal/message"
	"github.com/stonecool/livemusic-go/internal/task"
	"log"
	"sync"
	"time"
)

var (
	scheduler *Scheduler
	once      sync.Once
)

type Scheduler struct {
	cron   *cron.Cron
	jobMap map[int]cron.EntryID
	mu     sync.RWMutex
}

func GetScheduler() *Scheduler {
	once.Do(func() {
		scheduler = &Scheduler{
			cron:   cron.New(cron.WithSeconds()),
			jobMap: make(map[int]cron.EntryID),
		}
		scheduler.Start()
	})
	return scheduler
}

func (s *Scheduler) AddTask(task task.ITask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if oldEntryID, exists := s.jobMap[task.GetID()]; exists {
		s.cron.Remove(oldEntryID)
	}

	entryID, err := s.cron.AddFunc(task.GetCronSpec(), func() {
		s.executeTask(task)
	})
	if err != nil {
		return err
	}

	s.jobMap[task.GetID()] = entryID
	return nil
}

func (s *Scheduler) Start() {
	// 启动时加载所有任务
	tasks, err := task.GetAllCrawlTasks()
	if err != nil {
		log.Printf("Failed to load crawl tasks: %v", err)
		return
	}

	for _, task := range tasks {
		if err := s.AddTask(task); err != nil {
			log.Printf("Failed to add task %d: %v", task.GetID(), err)
		}
	}

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) executeTask(task task.ITask) error {
	const (
		maxRetries = 3 // 最大重试次数
		retryDelay = 5 // 重试间隔(秒)
	)

	msg := message.NewAsyncMessageWithCmd(message.AccountCmd_Crawl, task)
	var lastErr error
	for retry := 0; retry < maxRetries; retry++ {
		// 如果不是第一次尝试,等待一段时间
		if retry > 0 {
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}

		// 尝试分发任务
		if err := pool.GetPool().DispatchTask(task.GetCategory(), msg); err == nil {
			return nil
		} else {
			lastErr = err
			log.Printf("Task %d dispatch failed (attempt %d/%d): %v",
				task.GetID(), retry+1, maxRetries, err)
		}
	}

	return fmt.Errorf("dispatch task failed after %d attempts: %v",
		maxRetries, lastErr)
}
