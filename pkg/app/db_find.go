package app

import (
	"fmt"
	"gorm.io/gorm"
	"lls_api/common"
	"runtime"
)

// ============ 查询执行方法 ============

func (d *Database) ByID(id common.ID, dest any) error {
	err := d.db().First(dest, "id = ?", id).Error
	if err != nil {
		// 获取调用者信息（跳过1层：当前函数）
		if _, file, line, ok := runtime.Caller(1); ok {
			// 打印或包装错误，带上调用位置
			return fmt.Errorf("%s:%d: %w", file, line, err)
		}
	}
	return err
}

func (d *Database) Find(dest any, conds ...any) error {
	return d.db().Find(dest, conds...).Error
}

func (d *Database) First(dest any, conds ...any) error {
	err := d.db().First(dest, conds...).Error
	if err != nil {
		// 获取调用者信息（跳过1层：当前函数）
		if _, file, line, ok := runtime.Caller(1); ok {
			// 打印或包装错误，带上调用位置
			return fmt.Errorf("%s:%d: %w", file, line, err)
		}
	}
	return err
}

func (d *Database) Take(dest any, conds ...any) error {
	err := d.db().Take(dest, conds...).Error
	if err != nil {
		// 获取调用者信息（跳过1层：当前函数）
		if _, file, line, ok := runtime.Caller(1); ok {
			// 打印或包装错误，带上调用位置
			return fmt.Errorf("%s:%d: %w", file, line, err)
		}
	}
	return err
}

func (d *Database) Last(dest any, conds ...any) error {
	return d.gdb.Last(dest, conds...).Error
}

func (d *Database) Count(count *int64) error {
	return d.gdb.Count(count).Error
}

func (d *Database) Pluck(column string, dest any) error {
	return d.gdb.Pluck(column, dest).Error
}

func (d *Database) Scan(dest any) error {
	return d.gdb.Scan(dest).Error
}

// ============ 查询条件构造方法 ============

func (d *Database) Model(value any) *Database {
	d.gdb = d.gdb.Model(value)
	return d
}

func (d *Database) Where(query any, args ...any) *Database {
	d.gdb = d.gdb.Where(query, args...)
	return d
}

func (d *Database) Not(query any, args ...any) *Database {
	d.gdb = d.gdb.Not(query, args...)
	return d
}

func (d *Database) Or(query any, args ...any) *Database {
	d.gdb = d.gdb.Or(query, args...)
	return d
}

func (d *Database) Preload(query string, args ...any) *Database {
	d.gdb = d.gdb.Preload(query, args...)
	return d
}

func (d *Database) Joins(query string, args ...any) *Database {
	d.gdb = d.gdb.Joins(query, args...)
	return d
}

func (d *Database) InnerJoins(query string, args ...any) *Database {
	d.gdb = d.gdb.InnerJoins(query, args...)
	return d
}

func (d *Database) Group(name string) *Database {
	d.gdb = d.gdb.Group(name)
	return d
}

func (d *Database) Having(query any, args ...any) *Database {
	d.gdb = d.gdb.Having(query, args...)
	return d
}

func (d *Database) Order(value any) *Database {
	d.gdb = d.gdb.Order(value)
	return d
}

func (d *Database) Limit(limit int) *Database {
	d.gdb = d.gdb.Limit(limit)
	return d
}

func (d *Database) Offset(offset int) *Database {
	d.gdb = d.gdb.Offset(offset)
	return d
}

func (d *Database) Distinct(args ...any) *Database {
	d.gdb = d.gdb.Distinct(args...)
	return d
}

func (d *Database) Omit(columns ...string) *Database {
	d.gdb = d.gdb.Omit(columns...)
	return d
}

func (d *Database) Select(query any, args ...any) *Database {
	d.gdb = d.gdb.Select(query, args...)
	return d
}

// Session 用于设置查询选项（如 DryRun, NewDB 等）
func (d *Database) Session(config *gorm.Session) *Database {
	d.gdb = d.gdb.Session(config)
	return d
}

// Table 指定表名（用于子查询等）
func (d *Database) Table(name string) *Database {
	d.gdb = d.gdb.Table(name)
	return d
}

type ReadOnlyScopeFunc func(database *Database) *Database

// 内部转换函数：将 ReadOnlyScopeFunc 转为 GORM 的 ScopeFunc
func (d *Database) toGormScope(f ReadOnlyScopeFunc) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 创建一个临时 ReadOnlyQuery，包装 db
		tempQ := &Database{gdb: db}
		// 调用用户函数，但用户只能操作 ReadOnlyQuery
		modifiedQ := f(tempQ)
		// 返回修改后的 db
		return modifiedQ.gdb
	}
}

func (d *Database) Scopes(funcs ...ReadOnlyScopeFunc) *Database {
	// 转换为 GORM 的 ScopeFunc
	gormFuncs := make([]func(*gorm.DB) *gorm.DB, len(funcs))
	for i, f := range funcs {
		gormFuncs[i] = d.toGormScope(f)
	}
	d.gdb = d.gdb.Scopes(gormFuncs...)
	return d
}
