package service

import (
	"context"

	usercmd "github.com/arvinpaundra/cent/user/application/command/user"
	"github.com/arvinpaundra/cent/user/domain/user/repository"
)

type UpdateUserBalance struct {
	userReader repository.UserReader
	userWriter repository.UserWriter
	uow        repository.UnitOfWork
}

func NewUpdateUserBalance(
	userReader repository.UserReader,
	userWriter repository.UserWriter,
	uow repository.UnitOfWork,
) UpdateUserBalance {
	return UpdateUserBalance{
		userReader: userReader,
		userWriter: userWriter,
		uow:        uow,
	}
}

func (s UpdateUserBalance) Exec(ctx context.Context, command usercmd.UpdateUserBalance) error {
	tx, err := s.uow.Begin()
	if err != nil {
		return err
	}

	user, err := tx.UserReader().FindByIdForUpdate(ctx, command.UserId)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	user.UpdateBalance(command.Amount)

	user.MarkToBeUpdated()

	err = tx.UserWriter().Save(ctx, user)
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
