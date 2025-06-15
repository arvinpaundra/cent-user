package repository

type UnitOfWork interface {
	Begin() (UnitOfWorkProcessor, error)
}

type UnitOfWorkProcessor interface {
	UserReader() UserReader
	UserWriter() UserWriter

	Rollback() error
	Commit() error
}
