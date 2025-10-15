package middleware

import (
	"github.com/gin-gonic/gin"
	"lls_api/pkg/log"
	"runtime/debug"
)

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				lc := log.DefaultContext()
				if err, ok := err.(error); ok {
					lc.ErrorErr(err)
				}
				lc.Error(string(debug.Stack()))
				// 兜底处理：确保不会重复写入响应
				if ctx.Writer.Written() {
					// 响应已经写入，不修改状态码，只记录日志
					ctx.Abort()
				} else if ctx.IsAborted() {
					// 已经被 abort，不处理
				} else {
					// 未写入响应，返回 500
					ctx.AbortWithStatusJSON(500, gin.H{
						"code": 500,
						"msg":  "服务器内部错误.",
					})
				}
			}
		}()
		ctx.Next()
	}
}
