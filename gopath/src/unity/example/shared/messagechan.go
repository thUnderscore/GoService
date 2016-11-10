package shared

//MessageChan join channel and MessageHandler and implements MessageSender interface 
type MessageChan struct {
	chn chan *Message
	*MessageHandler
}

//NewMessageChan creates MessageChan. Define cap if you need buffered channel
func NewMessageChan(cap int) *MessageChan {
	return &MessageChan{
		chn:            make(chan *Message, cap),
		MessageHandler: NewMessageHandler()}
}

//Send adds message to channel and doesn't wait for result. BUt it still can block goroutine if chan buffer is full
func (mchn *MessageChan) Send(code MessageCode, data interface{}) {
	m := newMessage(code, data, false)
	mchn.chn <- m
}

//SendSync adds message to channel and wait for handling result.
func (mchn *MessageChan) SendSync(code MessageCode, data interface{}) interface{} {
	m := newMessage(code, data, true)
	mchn.chn <- m
	return m.wait()
}
