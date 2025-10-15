package app

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lls_api/common"
	"lls_api/pkg/config"
	"lls_api/pkg/log"
	"runtime"
)

type Database struct {
	gdb     *gorm.DB
	log     *log.LoggerContext // 后续可以替换成接口
	isDebug bool
}

func newDatabase(log *log.LoggerContext, isDebug bool) *Database {
	return &Database{gdb: gdb, log: log, isDebug: isDebug}
}

func newDatabaseWithTx(tx *gorm.DB, log *log.LoggerContext, isDebug bool) *Database {
	return &Database{gdb: tx, log: log, isDebug: isDebug}
}

func (d *Database) db() *gorm.DB {
	if d.isDebug {
		return d.gdb.Debug()
	}
	return d.gdb
}

func (d *Database) Log() *log.LoggerContext {
	return d.log
}

var gdb *gorm.DB

func InitDB() {
	var err error
	gdb, err = gorm.Open(mysql.Open(config.C.DB.Dsn))
	if err != nil {
		panic(err)
	}
	if config.C.DB.Debug {
		gdb = gdb.Debug()
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(config.C.DB.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.C.DB.MaxOpenConn)
	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}
}

type Moduler interface {
	TableName() string
	GetID() common.ID
}

type Result struct {
	RowsAffected int64
	Error        error
}

func wrapCaller(err error, skip ...int) error {
	if err == nil {
		return nil
	}
	sk := 2
	if len(skip) > 0 {
		sk += skip[0]
	}
	// 获取调用者信息（跳过1层：当前函数）
	if _, file, line, ok := runtime.Caller(sk); ok {
		// 打印或包装错误，带上调用位置
		return fmt.Errorf("%s:%d: %w", file, line, err)
	}
	return err
}

func (d *Database) AutoMigrateModels(models []any) error {
	for _, m := range models {
		if err := d.gdb.AutoMigrate(m); err != nil {
			return err
		}
	}

	return nil
}
