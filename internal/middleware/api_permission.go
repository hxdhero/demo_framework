package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"lls_api/pkg/app"
	"net/http"
)

type CachePermission struct {
	Actions []string `json:"actions"`
}

// ApiPermission 校验http接口请求权限的中间件,应用此中间件后如果没有权限返回403
func ApiPermission() gin.HandlerFunc {
	return app.WrapHandler(func(ctx *app.Context) {
		reqPermission := ctx.Rest().Gin().FullPath() + ":" + ctx.Rest().Gin().Request.Method
		ctx.Log().Info(reqPermission)
		permissionKey := fmt.Sprintf("actions|Bearer %s", ctx.Session().AuthToken)
		var cachePermission CachePermission
		bs, err := ctx.Redis().Get(ctx.Context(), permissionKey).Bytes()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				ctx.Log().Infof("验证权限错误, key没有找到")
			} else {
				ctx.Log().Infof("验证权限错误,获取权限信息错误: %e", err)
			}
			ctx.Rest().JSON(http.StatusForbidden, map[string]interface{}{"detail": "您没有执行该操作的权限。"})
			ctx.Rest().Abort()
			return
		}
		if err := json.Unmarshal(bs, &cachePermission); err != nil {
			ctx.Log().Infof("验证权限错误,获取权限信息错误: %e", err)
			ctx.Rest().JSON(http.StatusForbidden, map[string]interface{}{"detail": "您没有执行该操作的权限。"})
			ctx.Rest().Abort()
			return
		}
		hasPermission := false
		for _, e := range cachePermission.Actions {
			if e == reqPermission {
				hasPermission = true
			}
		}
		if !hasPermission {
			ctx.Log().Infof("[您没有执行该操作的权限] 请为 actions 配置: \"%s\"", reqPermission)
			ctx.Rest().JSON(http.StatusForbidden, map[string]interface{}{"detail": "您没有执行该操作的权限。"})
			ctx.Rest().Abort()
			return
		}

	})
}

// AuthAndPermission 带鉴权和权限的中间件
func AuthAndPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		Auth()(c)
		if c.IsAborted() {
			return
		}
		ApiPermission()(c)
		if c.IsAborted() {
			return
		}
	}
}
