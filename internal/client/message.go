package client

import "github.com/stonecool/livemusic-go/internal"

type AsyncMessage struct {
	*internal.Message                  // 嵌入原有的通信 Message
	Result            chan interface{} // 异步结果通道
}

func NewAsyncMessage(msg *internal.Message) *AsyncMessage {
	return &AsyncMessage{
		Message: msg,
		Result:  make(chan interface{}),
	}
}
