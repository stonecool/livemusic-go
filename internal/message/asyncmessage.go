package message

type AsyncMessage struct {
	*Message
	Result chan *Message
}

func NewAsyncMessage(msg *Message) *AsyncMessage {
	return &AsyncMessage{
		Message: msg,
		Result:  make(chan *Message, 1),
	}
}
