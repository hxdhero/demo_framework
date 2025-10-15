package middleware

import (
	"github.com/gin-gonic/gin"
	"lls_api/pkg/app"
	"slices"
)

var permissionRule map[string][]string

func init() {
	permissionRule = map[string][]string{
		// "/1/g/1.0/example/:id/":      nil,
		// "/1/g/1.0/example/:id/edit/": {"labor_cdl", "labor_lls"},
	}
}

// DefaultPermission 默认的权限中间件
//
// 配置:
//
//	permissionRule = map[string][]string{
//	"/1/g/1.0/example/:id/":      nil,
//	"/1/g/1.0/example/:id/edit/": {"labor_cdl", "labor_lls"},
//	}
//
// 效果:
//
//	url: "/1/g/1.0/example/:id/"  不会校验任何权限
//	url: "/1/g/1.0/example/:id/edit/" 如果jwt中的userType是labor_cdl或者labor_lls 校验就可以通过
//	url: "/1/g/1.0/example/:id/other/" 会被校验菜单权限(redis)
func DefaultPermission(ps ...map[string][]string) gin.HandlerFunc {
	if len(ps) > 0 {
		permissionRule = ps[0]
	}
	return app.WrapHandler(func(ctx *app.Context) {
		url := ctx.Rest().Gin().FullPath()
		vals, ok := permissionRule[url]
		if ok {
			if len(vals) == 0 {
				ctx.Rest().Gin().Next()
				return
			}
			if slices.Contains(vals, ctx.Session().Auth.UserType) {
				ctx.Rest().Gin().Next()
				return
			}
		}
		AuthAndPermission()(ctx.Rest().Gin())
	})
}
