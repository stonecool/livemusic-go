package message

import "github.com/stonecool/livemusic-go/internal/task"

type AsyncMessage struct {
	*Message
	Result            chan interface{}
	*task.Task
}

func NewAsyncMessage(msg *Message, task *task.Task) *AsyncMessage {
	return &AsyncMessage{
		Message: msg,
		Result:  make(chan interface{}),
		Task:task,
	}
}
