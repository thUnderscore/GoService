package shared

import "testing"

var _ MessageSender = new(MessageChan) //should implement MessageSender
func TestMessageChan(t *testing.T) {

	mchn := NewMessageChan(0)
	val := 0
	mchn.SetHandler(0, func(m *Message) {
		val = val + m.Data.(int)
		m.Data = val
	})
	mchn.SetHandler(1, func(m *Message) {
		val = val * m.Data.(int)
		m.Data = val
	})
	j := NewJobChan(func(j *JobChan) {
		for {
			var m *Message
			select {
			case <-j.ExitChn:
				return
			case m = <-mchn.Chn:
				mchn.Handle(m)
			}
		}
	})
	j.Start(false)

	mchn.Send(0, 2)
	res := mchn.SendSync(1, 3).(int)
	if res != val {
		t.Error("res = ", res, "val = ", val)
	}
	if res != 6 {
		t.Error("res = ", res, " should be", 6)
	}

	j.Stop(true)
}
