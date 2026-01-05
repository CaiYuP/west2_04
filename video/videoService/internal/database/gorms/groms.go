package gorms

import (
	"context"
	"gorm.io/gorm"
)

var _db *gorm.DB

//func init() {
//	if config.C.Dbc.Separation {
//		username := config.C.Dbc.Master.Username
//		password := config.C.Dbc.Master.Password
//		host := config.C.Dbc.Master.Host
//		port := config.C.Dbc.Master.Port
//		dbName := config.C.Dbc.Master.Db
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
//		var err error
//		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//			Logger: logger.Default.LogMode(logger.Info),
//		})
//		if err != nil {
//			panic("数据库连接失败:" + err.Error())
//		}
//		replicas := []gorm.Dialector{}
//		for _, v := range config.C.Dbc.Slave {
//			username := v.Username //账号
//			password := v.Password //密码
//			host := v.Host         //数据库地址，可以是Ip或者域名
//			port := v.Port         //数据库端口
//			Dbname := v.Db         //数据库名
//			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
//			cfg := mysql.Config{
//				DSN: dsn,
//			}
//			replicas = append(replicas, mysql.New(cfg))
//		}
//		_db.Use(dbresolver.Register(dbresolver.Config{
//			//主库
//			Sources: []gorm.Dialector{mysql.New(mysql.Config{
//				DSN: dsn,
//			})},
//			//从库
//			Replicas: replicas,
//			Policy:   dbresolver.RandomPolicy{},
//		}).
//			SetMaxIdleConns(10).
//			SetMaxOpenConns(200))
//	} else {
//		username := config.C.Mc.Username
//		password := config.C.Mc.Password
//		host := config.C.Mc.Host
//		port := config.C.Mc.Port
//		dbName := config.C.Mc.Db
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
//		var err error
//		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//			Logger: logger.Default.LogMode(logger.Info),
//		})
//		if err != nil {
//			panic("数据库连接失败:" + err.Error())
//		}
//
//	}
//
//}
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
