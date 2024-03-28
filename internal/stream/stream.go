package stream

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func NewNats(natsUrl, clientId string) (stan.Conn, error) {
	natsConn, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}

	stanConn, err := stan.Connect("test-cluster", clientId, stan.NatsConn(natsConn), stan.ConnectWait(time.Second * 15))
	if err != nil {
		return nil, err
	}

	return stanConn, nil
}