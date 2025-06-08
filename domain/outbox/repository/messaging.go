package repository

import "context"

type Messaging interface {
	Publish(ctx context.Context, topic string, payload []byte) error
}
