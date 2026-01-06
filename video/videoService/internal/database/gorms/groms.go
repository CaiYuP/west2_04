package gorms

import (
	"context"
	"gorm.io/gorm"
)

var _db *gorm.DB

func GetDB() *gorm.DB {
	return _db
}
func SetDB(db *gorm.DB) {
	_db = db
}

type GormConn struct {
	db *gorm.DB
	//事务连接
	tx *gorm.DB
}

func (g *GormConn) Rollback() error {
	return g.tx.Rollback().Error
}

func (g *GormConn) Commit() error {
	return g.tx.Commit().Error
}
func (g *GormConn) Begin() {
	g.tx = g.db.Begin()
}
func New() *GormConn {
	return &GormConn{
		db: GetDB(),
	}
}
func NewTran() *GormConn {
	return &GormConn{
		db: GetDB(),
		tx: GetDB(),
	}
}
func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.tx.WithContext(ctx)
}
