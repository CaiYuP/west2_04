package database

type DbConn interface {
	Begin()
	Rollback() error
	Commit() error
}
