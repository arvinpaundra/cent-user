package service

import (
	"context"
	"encoding/json"

	authcmd "github.com/arvinpaundra/cent/user/application/command/auth"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type Register struct {
	userReader   repository.UserReader
	userWriter   repository.UserWriter
	outboxWriter repository.OutboxWriter
	unitOfWork   repository.UnitOfWork
}

func NewRegister(
	userReader repository.UserReader,
	userWriter repository.UserWriter,
	outboxWriter repository.OutboxWriter,
	unitOfWork repository.UnitOfWork,
) Register {
	return Register{
		userReader:   userReader,
		userWriter:   userWriter,
		outboxWriter: outboxWriter,
		unitOfWork:   unitOfWork,
	}
}

func (s Register) Exec(ctx context.Context, payload authcmd.Register) error {
	isExist, err := s.userReader.IsEmailExist(ctx, payload.Email)
	if err != nil {
		return err
	}

	if isExist {
		return constant.ErrEmailAlreadyTaken
	}

	user := entity.User{
		Email:    payload.Email,
		Fullname: payload.Fullname,
	}

	err = user.GeneratePassword(payload.Password)
	if err != nil {
		return err
	}

	err = user.GenerateKey()
	if err != nil {
		return err
	}

	tx, err := s.unitOfWork.Begin()
	if err != nil {
		return err
	}

	err = tx.UserWriter().Save(ctx, &user)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	err = user.GenerateSlug()
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	user.MarkToBeUpdated()

	err = tx.UserWriter().Save(ctx, &user)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	outboxPayload := struct {
		UserId   int64   `json:"user_id"`
		UserSlug *string `json:"user_slug"`
	}{
		UserId:   user.ID,
		UserSlug: user.Slug,
	}

	payloadBytes, err := json.Marshal(outboxPayload)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	outbox := entity.Outbox{
		Event:   constant.OutboxEventUserRegistered,
		Status:  constant.OutboxStatusPending,
		Payload: payloadBytes,
	}

	err = tx.OutboxWriter().Save(ctx, &outbox)
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
