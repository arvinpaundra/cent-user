package repository

type UnitOfWork interface {
	Begin() (UnitOfWorkProcessor, error)
}

type UnitOfWorkProcessor interface {
	OutboxWriter() OutboxWriter

	Rollback() error
	Commit() error
}
