package utils

import "log"

type Notification struct {
	Type string
	Data string
}

type Notifier struct {
	ch chan Notification
}

func NewNotifier(buffer int) *Notifier {
	return &Notifier{
		ch: make(chan Notification, buffer),
	}
}

func (n *Notifier) Start() {
	go func() {
		for msg := range n.ch {
			log.Printf("[NOTIFY] %s: %s\n", msg.Type, msg.Data)
		}
	}()
}

func (n *Notifier) Send(msg Notification) {
	select {
	case n.ch <- msg:
	default:
	}
}
