package user

import (
	"github.com/arvinpaundra/cent/user/domain/user/repository"
	"gorm.io/gorm"
)

var _ repository.UnitOfWork = (*UnitOfWork)(nil)
var _ repository.UnitOfWorkProcessor = (*UnitOfWorkProcessor)(nil)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return UnitOfWork{db: db}
}

func (r UnitOfWork) Begin() (repository.UnitOfWorkProcessor, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return UnitOfWorkProcessor{tx: tx}, nil
}

type UnitOfWorkProcessor struct {
	tx *gorm.DB
}

func (r UnitOfWorkProcessor) UserReader() repository.UserReader {
	return NewUserReaderRepository(r.tx)
}

func (r UnitOfWorkProcessor) UserWriter() repository.UserWriter {
	return NewUserWriterRepository(r.tx)
}

func (r UnitOfWorkProcessor) Rollback() error {
	return r.tx.Rollback().Error
}

func (r UnitOfWorkProcessor) Commit() error {
	return r.tx.Commit().Error
}
