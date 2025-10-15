package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lls_api/pkg/log"
	"lls_api/pkg/rdb"
)

const (
	TxKey = "app_context_tx"
)

type Context struct {
	ctx         context.Context
	database    *Database
	rdb         *rdb.RedisClient
	rest        *Rest
	logContext  *log.LoggerContext
	sessionInfo *SessionInfo
}

func NewAppContext(ctx context.Context, ginContext *gin.Context, logContext *log.LoggerContext, sessionInfo *SessionInfo) *Context {
	return &Context{
		ctx:         ctx,
		database:    newDatabase(logContext, sessionInfo.IsDebug),
		rest:        NewRest(ginContext, logContext),
		logContext:  logContext,
		sessionInfo: sessionInfo,
	}
}

func ContextWithTest() *Context {
	return NewAppContext(context.Background(), nil, log.DefaultContext(), &SessionInfo{})
}

// DB 获取封装的db对象,如果当前外层代码在事务中就会加入事务
func (c *Context) DB() *Database {
	if c.ctx != nil {
		if tx, ok := c.ctx.Value(TxKey).(*gorm.DB); ok {
			return newDatabaseWithTx(tx, c.logContext, c.sessionInfo.IsDebug)
		}
	}
	if c.sessionInfo.IsDebug {
		c.database.isDebug = true
	}
	return c.database
}

func (c *Context) Rest() *Rest {
	return c.rest
}

// DBWithTransaction 开启事务
// 如果外层代码已经开启了事务 ctx.DBWithTransaction 会加入到外层的事务
// 开启事务后 调用ctx.DB() 也会加入到当前事务或者外层事务
func (c *Context) DBWithTransaction(fn func(ctx *Context) error) error {
	if c.ctx != nil {
		if _, ok := c.ctx.Value(TxKey).(*gorm.DB); ok {
			return fn(c) // ← 不 Commit/Rollback，交给外层
		}
	}

	tx := c.database.gdb.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 创建带事务上下文的新 AppInfo（不修改原对象！）
	txCtx := context.WithValue(c.ctx, TxKey, tx)
	txContext := NewAppContext(txCtx, c.rest.gin, c.logContext, c.sessionInfo)
	err := fn(txContext)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Redis 返回封装的redis实例
func (c *Context) Redis() *rdb.RedisClient {
	return c.rdb
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Session() *SessionInfo {
	return c.sessionInfo
}

func (c *Context) Log() *log.LoggerContext {
	if c.logContext == nil {
		c.logContext = log.DefaultContext()
	}
	return c.logContext
}

type IOriginalError interface {
	OriginalErr() error
}
