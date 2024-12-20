package task

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"go.uber.org/zap"
)

type Queue struct {
	taskChan chan ITask
	done     chan struct{}
}

var DefaultQueue *Queue

func init() {
	DefaultQueue = NewQueue(100)
}

func NewQueue(size int) *Queue {
	return &Queue{
		taskChan: make(chan ITask, size),
		done:     make(chan struct{}),
	}
}

func (q *Queue) Stop() {
	close(q.done)
}

func (q *Queue) PushTask(task ITask) error {
	select {
	case q.taskChan <- task:
		internal.Logger.Info("Task pushed to queue",
			zap.String("category", task.GetCategory()))
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (q *Queue) PopTaskByCategory(category string) (ITask, error) {
	select {
	case task := <-q.taskChan:
		if task.GetCategory() == category {
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
