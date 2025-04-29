package events

import (
	"log"

	"github.com/nats-io/nats.go"
)

type Msg struct {
	Subject string
	Data    []byte
}

var nc *nats.Conn

func Connect(url string) {
	var err error
	nc, err = nats.Connect(url)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
}

func Publish(subject string, data []byte) error {
	return nc.Publish(subject, data)
}

type EventHandler func(*Msg)

func Subscribe(subject string, handler EventHandler) (*nats.Subscription, error) {
	return nc.Subscribe(subject, func(msg *nats.Msg) {
		eventMsg := &Msg{
			Subject: msg.Subject,
			Data:    msg.Data,
		}
		handler(eventMsg)
	})
}
