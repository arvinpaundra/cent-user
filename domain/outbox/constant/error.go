package constant

import "errors"

var (
	ErrOutboxNotFound       = errors.New("outbox not found")
	ErrFailedToPublishEvent = errors.New("failed to publish outbox event")
)
