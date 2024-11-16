package internal

import (
	"github.com/robfig/cron/v3"
	"log"
	"sync"
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

func (s *Scheduler) AddTask(task *CrawlTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果任务已存在,先移除
	if oldEntryID, exists := s.jobMap[task.ID]; exists {
		s.cron.Remove(oldEntryID)
	}

	// 添加新任务
	entryID, err := s.cron.AddFunc(task.cronSpec, func() {
		s.executeTask(task)
	})
	if err != nil {
		return err
	}

	s.jobMap[task.ID] = entryID
	return nil
}

func (s *Scheduler) Start() {
	// 启动时加载所有任务
	tasks, err := GetAllCrawlTasks()
	if err != nil {
		log.Printf("Failed to load crawl tasks: %v", err)
		return
	}

	for _, task := range tasks {
		if err := s.AddTask(task); err != nil {
			log.Printf("Failed to add task %d: %v", task.ID, err)
		}
	}

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) executeTask(task *CrawlTask) {
	if err := task.Execute(); err != nil {
		log.Printf("Failed to execute task %d: %v", task.ID, err)
		return
	}
}
