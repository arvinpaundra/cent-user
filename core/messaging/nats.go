package messaging

import (
	"log"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	conn *nats.Conn
}

func NewNats(url string) *Nats {
	conn, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("failed to connect to nats: %s", err.Error())
	}

	return &Nats{
		conn: conn,
	}
}

func (n *Nats) GetConnection() *nats.Conn {
	return n.conn
}

func (n *Nats) Close() error {
	return n.conn.Drain()
}
