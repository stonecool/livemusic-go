package client

type AsyncMessage struct {
	*Message
	Result            chan interface{}
}

func NewAsyncMessage(msg *Message) *AsyncMessage {
	return &AsyncMessage{
		Message: msg,
		Result:  make(chan interface{}),
	}
}
