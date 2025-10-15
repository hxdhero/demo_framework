package middleware

import (
	"github.com/gin-gonic/gin"
	"lls_api/pkg/app"
	"strings"
	"time"
)

func Log() gin.HandlerFunc {
	return app.WrapHandler(func(ctx *app.Context) {
		// 记录请求开始时间
		startTime := time.Now()
		urlPath := ctx.Rest().Gin().Request.URL.Path
		ctx.Log().Infof("[lls-request] method:%s; path:%s; %s", ctx.Rest().Gin().Request.Method, urlPath, strings.Join(strings.Split(urlPath, "/"), "_"))
		ctx.Rest().Gin().Next()
		// 计算处理时间
		duration := time.Since(startTime)
		ctx.Log().Infof("[lls-response] status_code:%d; duration:%v", ctx.Rest().Gin().Writer.Status(), duration)
	})
}
