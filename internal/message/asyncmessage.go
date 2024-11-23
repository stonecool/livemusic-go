package message

import "github.com/stonecool/livemusic-go/internal/task"

type AsyncMessage struct {
	*Message
	Result chan interface{}
	*task.Task
}

func NewAsyncMessageWithMsg(msg *Message, task *task.Task) *AsyncMessage {
	return &AsyncMessage{
		Message: msg,
		Result:  make(chan interface{}),
		Task:    task,
	}
}

func NewAsyncMessageWithCmd(cmd CrawlCmd, task *task.Task) *AsyncMessage {
	return &AsyncMessage{
		Message: &Message{
			Cmd: cmd,
		},
		Result: make(chan interface{}),
		Task:   task,
	}
}
