package internal

type AsyncMessage struct {
	*Message                  // 嵌入原有的通信 Message
	Result   chan interface{} // 异步结果通道
}

func NewAsyncMessage(msg *Message) *AsyncMessage {
	return &AsyncMessage{
		Message: msg,
		Result:  make(chan interface{}),
	}
}
