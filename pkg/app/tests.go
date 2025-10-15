package app

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ErrTestRollback = errors.New("test: rollback transaction intentionally")

// WithTestTransaction 用于测试的事务,最终会被rollback
func WithTestTransaction(t *testing.T, ctx *Context, testFunc func(ctx *Context) error) {
	// t.Helper()
	err := ctx.DBWithTransaction(func(ctx *Context) error {
		if innerErr := testFunc(ctx); innerErr != nil {
			return innerErr // 真实错误，不包装
		}
		return ErrTestRollback // 总是回滚
	})
	assert.ErrorIs(t, err, ErrTestRollback)
}

// WithTestMTransaction 用于测试的事务,最终会被rollback
func WithTestMTransaction(m *testing.M, ctx *Context, testFunc func(ctx *Context) error) {
	// t.Helper()
	err := ctx.DBWithTransaction(func(ctx *Context) error {
		if innerErr := testFunc(ctx); innerErr != nil {
			return innerErr // 真实错误，不包装
		}
		return ErrTestRollback // 总是回滚
	})
	if !errors.Is(err, ErrTestRollback) {
		panic(err)
	}
}
