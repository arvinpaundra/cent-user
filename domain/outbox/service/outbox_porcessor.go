package service

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/outbox/constant"
	"github.com/arvinpaundra/cent/user/domain/outbox/repository"
)

type OutboxProcessor struct {
	outboxReader repository.OutboxReader
	outboxWriter repository.OutboxWriter
	unitOfWork   repository.UnitOfWork
	messaging    repository.Messaging
}

func NewOutboxProcessor(
	outboxReader repository.OutboxReader,
	outboxWriter repository.OutboxWriter,
	unitOfWork repository.UnitOfWork,
	messaging repository.Messaging,
) OutboxProcessor {
	return OutboxProcessor{
		outboxReader: outboxReader,
		outboxWriter: outboxWriter,
		unitOfWork:   unitOfWork,
		messaging:    messaging,
	}
}

func (s OutboxProcessor) Exec(ctx context.Context) error {
	outbox, err := s.outboxReader.FindUnprocessed(ctx)
	if err != nil {
		return err
	}

	tx, err := s.unitOfWork.Begin()
	if err != nil {
		return err
	}

	outbox.MarkToBeUpdated()

	outbox.SetStatus(constant.OutboxStatusProcessing)

	err = tx.OutboxWriter().Save(ctx, outbox)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	topic := s.topic(outbox.Event)

	err = s.messaging.Publish(ctx, topic, outbox.Payload)
	if err != nil {
		outbox.SetStatus(constant.OutboxStatusFailed)
		outbox.SetError(err.Error())

		err = tx.OutboxWriter().Save(ctx, outbox)
		if err != nil {
			if uowErr := tx.Rollback(); uowErr != nil {
				return uowErr
			}

			return err
		}

		if uowErr := tx.Commit(); uowErr != nil {
			return uowErr
		}

		return constant.ErrFailedToPublishEvent
	}

	outbox.SetStatus(constant.OutboxStatusPublished)
	outbox.SetPublishedAt()

	err = tx.OutboxWriter().Save(ctx, outbox)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	if uowErr := tx.Commit(); uowErr != nil {
		return uowErr
	}

	return nil
}

func (s OutboxProcessor) topic(event constant.OutboxEvent) string {
	switch event {
	case constant.OutboxEventUserRegistered:
		return constant.EventUserCreated
	default:
		return ""
	}
}
