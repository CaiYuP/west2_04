package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"time"
	"videoService/internal/database/gorms"
	"west2-video/common/logs"
)

var _db *gorm.DB

// ReConnMysql 连接MySQL数据库（支持读写分离）
func (c *Config) ReConnMysql() {
	if c.Dbc.Separation {
		// 读写分离配置
		username := c.Dbc.Master.Username // 账号
		password := c.Dbc.Master.Password // 密码
		host := c.Dbc.Master.Host         // 数据库地址
		port := c.Dbc.Master.Port         // 数据库端口
		Dbname := c.Dbc.Master.Db         // 数据库名

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
			username, password, host, port, Dbname)

		var err error
		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if err != nil {
			logs.LG.Error("连接主数据库失败, error=" + err.Error())
			return
		}

		// 配置从库连接
		replicas := []gorm.Dialector{}
		for _, v := range c.Dbc.Slave {
			slaveDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
				v.Username, v.Password, v.Host, v.Port, v.Db)

			cfg := mysql.Config{
				DSN:                       slaveDsn,
				DefaultStringSize:         256,   // string 类型字段的默认长度
				DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前不支持
				DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前不支持
				DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前不支持
				SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
			}
			replicas = append(replicas, mysql.New(cfg))
		}

		// 配置读写分离插件
		_db.Use(dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: dsn,
			})},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{}, // 随机选择从库
			// 以下是连接池配置
		}).
			SetMaxIdleConns(c.getMaxIdleConns()).
			SetMaxOpenConns(c.getMaxOpenConns()).
			SetConnMaxLifetime(c.getConnMaxLifetime()))

		logs.LG.Info("MySQL读写分离模式连接成功")

	} else {
		// 单数据库配置
		username := c.Mc.Username // 账号
		password := c.Mc.Password // 密码
		host := c.Mc.Host         // 数据库地址
		port := c.Mc.Port         // 数据库端口
		Dbname := c.Mc.Db         // 数据库名

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
			username, password, host, port, Dbname)

		var err error
		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if err != nil {
			logs.LG.Error("连接MySQL数据库失败, error=" + err.Error())
			return
		}

		// 配置连接池
		sqlDB, err := _db.DB()
		if err != nil {
			logs.LG.Error("获取数据库连接池失败, error=" + err.Error())
			return
		}

		sqlDB.SetMaxIdleConns(c.getMaxIdleConns())
		sqlDB.SetMaxOpenConns(c.getMaxOpenConns())
		sqlDB.SetConnMaxLifetime(c.getConnMaxLifetime())

		logs.LG.Info("MySQL单机模式连接成功")
	}

	// 设置全局数据库连接
	gorms.SetDB(_db)

	// 测试连接
	c.testConnection()
}

// testConnection 测试数据库连接
func (c *Config) testConnection() {
	if _db == nil {
		logs.LG.Error("数据库连接为空")
		return
	}

	sqlDB, err := _db.DB()
	if err != nil {
		logs.LG.Error("获取数据库连接池失败, error=" + err.Error())
		return
	}

	if err := sqlDB.Ping(); err != nil {
		logs.LG.Error("数据库连接测试失败, error=" + err.Error())
		return
	}

	logs.LG.Info("数据库连接测试成功")
}

// getMaxIdleConns 获取最大空闲连接数
func (c *Config) getMaxIdleConns() int {
	if c.Dbc.MaxIdleConns > 0 {
		return c.Dbc.MaxIdleConns
	}
	return 10 // 默认值
}

// getMaxOpenConns 获取最大打开连接数
func (c *Config) getMaxOpenConns() int {
	if c.Dbc.MaxOpenConns > 0 {
		return c.Dbc.MaxOpenConns
	}
	return 200 // 默认值
}

// getConnMaxLifetime 获取连接最大生命周期
func (c *Config) getConnMaxLifetime() time.Duration {
	if c.Dbc.ConnMaxLifetime > 0 {
		return time.Duration(c.Dbc.ConnMaxLifetime) * time.Second
	}
	return time.Hour // 默认1小时
}

// GetDB 获取数据库实例（对外暴露）
func GetDB() *gorm.DB {
	if _db == nil {
		logs.LG.Warn("数据库连接尚未初始化，请先调用ReConnMysql方法")
	}
	return _db
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if _db != nil {
		sqlDB, err := _db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
