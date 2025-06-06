package constant

type OutboxEvent string

func (c OutboxEvent) String() string {
	return string(c)
}

const (
	OutboxEventUserRegistered OutboxEvent = "UserRegistered"
)

type OutboxStatus string

func (c OutboxStatus) String() string {
	return string(c)
}

const (
	OutboxStatusPending    OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusPublished  OutboxStatus = "published"
	OutboxStatusFailed     OutboxStatus = "failed"
)
