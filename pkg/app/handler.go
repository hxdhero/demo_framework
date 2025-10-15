package app

import (
	"github.com/gin-gonic/gin"
	"lls_api/common"
	"lls_api/pkg/config"
)

type HandlerFunc func(ctx *Context)

type SessionInfo struct {
	RequestID string
	Auth      common.Auth
	AuthErr   error
	AuthToken string
	Instance  common.UserInstance
	IsDebug   bool
}

// WrapHandler 把gin.Context 转为封装后的 app.Context
func WrapHandler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isTest {
			if val, ok := c.Get("app_context"); ok && isInit {
				ctx, ok := val.(*Context)
				if !ok {
					panic("app_context 类型错误.")
				}
				handler(ctx)
				return
			}
			appContext := ContextWithGin(c)
			appContext.ctx = testCtx.ctx
			isInit = true
			c.Set("app_context", appContext)
			handler(appContext)
		} else {
			// 如果已经存在app.Context 说明已经初始化过了直接返回. 因为后续中间件可能会多次调用WrapHandler方法
			if val, ok := c.Get("app_context"); ok {
				ctx, ok := val.(*Context)
				if !ok {
					panic("app_context 类型错误.")
				}
				handler(ctx)
				return
			}
			appContext := ContextWithGin(c)
			c.Set("app_context", appContext)
			handler(appContext)
		}

	}
}

var (
	isTest  bool
	isInit  bool
	testCtx *Context
)

func SetTestTransaction(ctx *Context) {
	if config.C.Env != config.EnvProd {
		isTest = true
		testCtx = ctx
	}
}
