package shared

//MessageHandler handler for Message that use Message.code to find corespondent function
type MessageHandler struct {
	handlers map[MessageCode]func(m *Message, isOn bool)
}

//NewMessageHandler creats new MessageHandler
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{handlers: make(map[MessageCode]func(m *Message, isOn bool))}
}

//SetHandler sets handler for code. If f is nil handler removes
func (h *MessageHandler) SetHandler(code MessageCode, f func(m *Message, isOn bool)) {
	if f == nil {
		delete(h.handlers, code)
	} else {
		h.handlers[code] = f
	}
}

//Handle handles Message
func (h *MessageHandler) Handle(msg *Message, isOn bool) {
	if handler, ok := h.handlers[msg.code]; ok {
		handler(msg, isOn)
	}
}
