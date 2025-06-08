package outbox

import (
	"context"
	"errors"

	"github.com/arvinpaundra/cent/user/domain/outbox/constant"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Messaging struct {
	conn *nats.Conn
}

func NewMessaging(conn *nats.Conn) Messaging {
	return Messaging{conn: conn}
}

func (r Messaging) Publish(ctx context.Context, topic string, payload []byte) error {
	js, err := jetstream.New(r.conn)
	if err != nil {
		return err
	}

	_, err = js.Stream(ctx, constant.StreamUser)
	if err != nil && !errors.Is(err, jetstream.ErrStreamNotFound) {
		return err
	}

	if errors.Is(err, jetstream.ErrStreamNotFound) {
		_, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:     constant.StreamUser,
			Subjects: []string{constant.EventUserCreated},
		})

		if err != nil {
			return err
		}
	}

	_, err = js.Publish(ctx, topic, payload)
	if err != nil {
		return err
	}

	return nil
}
