package message

import "github.com/stonecool/livemusic-go/internal/task"

type AsyncMessage struct {
	*Message
	Result chan interface{}
	task.ITask
}

func NewAsyncMessageWithMsg(msg *Message, task task.ITask) *AsyncMessage {
	return &AsyncMessage{
		Message: msg,
		Result:  make(chan interface{}),
		ITask:   task,
	}
}

func NewAsyncMessageWithCmd(cmd AccountCmd, task task.ITask) *AsyncMessage {
	return &AsyncMessage{
		Message: &Message{
			Cmd: cmd,
		},
		Result: make(chan interface{}),
		ITask:  task,
	}
}
