package repository

type UnitOfWork interface {
	Begin() (UnitOfWorkProcessor, error)
}

type UnitOfWorkProcessor interface {
	SessionWriter() SessionWriter
	UserWriter() UserWriter
	OutboxWriter() OutboxWriter

	Rollback() error
	Commit() error
}
