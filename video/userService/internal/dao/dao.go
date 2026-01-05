package dao

import (
	"errors"
	"userService/internal/database"
	"userService/internal/database/gorms"
	"west2-video/common/errs"
)

type TransactionDao struct {
	conn *gorms.GormConn
}

func (t *TransactionDao) Action(f func(conn database.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	var bErr *errs.BError
	if errors.Is(err, bErr) {
		errors.As(err, &bErr)
		if bErr != nil {
			t.conn.Rollback()
			return bErr
		} else {
			t.conn.Commit()
			return nil
		}
	}
	if err != nil {
		t.conn.Rollback()
		return err
	}
	return t.conn.Commit()
}

func NewTransactionDao() *TransactionDao {
	return &TransactionDao{
		conn: gorms.NewTran(),
	}
}
