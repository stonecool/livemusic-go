package message

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"go.uber.org/zap"
)

type CrawlTask struct {
	Category string
	Message  *AsyncMessage
}

type Queue struct {
	taskChan chan *CrawlTask
	done     chan struct{}
}

var DefaultQueue *Queue

func init() {
	DefaultQueue = NewQueue(100)
	go DefaultQueue.Start()
}

func NewQueue(size int) *Queue {
	return &Queue{
		taskChan: make(chan *CrawlTask, size),
		done:     make(chan struct{}),
	}
}

func (q *Queue) Start() {
	for {
		select {
		case <-q.done:
			return
		}
	}
}

func (q *Queue) Stop() {
	close(q.done)
}

func (q *Queue) PushTask(task *CrawlTask) error {
	select {
	case q.taskChan <- task:
		internal.Logger.Info("Task pushed to queue",
			zap.String("category", task.Category))
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (q *Queue) PopTaskByCategory(category string) (*CrawlTask, error) {
	select {
	case task := <-q.taskChan:
		if task.Category == category {
			return task, nil
		}
		// 如果类别不匹配，放回队列
		go func() {
			q.PushTask(task)
		}()
		return nil, fmt.Errorf("no matching task")
	default:
		return nil, fmt.Errorf("queue is empty")
	}
}
