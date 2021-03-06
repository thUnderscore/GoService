package shared

//MessageHandler handler for Message that use Message.code to find corespondent function
type MessageHandler struct {
	handlers map[MessageCode]func(m *Message)
}

//NewMessageHandler creats new MessageHandler
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{handlers: make(map[MessageCode]func(m *Message))}
}

//SetHandler sets handler associated with message code. If f is nil handler removes
func (h *MessageHandler) SetHandler(code MessageCode, f func(m *Message)) {
	if f == nil {
		delete(h.handlers, code)
	} else {
		h.handlers[code] = f
	}
}

//Handle handles Message
func (h *MessageHandler) Handle(m *Message) {
	if handler, ok := h.handlers[m.Code]; ok {
		m.Handle(handler)
	} else {
		m.free()
	}
}
