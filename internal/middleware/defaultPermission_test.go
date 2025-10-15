package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"lls_api/common"
	"lls_api/pkg"
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/util/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDefaultPermission(t *testing.T) {
	pkg.InitGlobal()

	// 初始化依赖
	pkg.InitGlobal()
	// 初始化appContext
	ctx := app.ContextWithTest()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(Recover())
	router.Use(CORS())
	router.Use(Log())
	router.Use(DefaultPermission(map[string][]string{
		"/1/g/1.0/example/:id/":      nil,
		"/1/g/1.0/example/:id/edit/": {"labor_cdl", "labor_lls"},
	}))
	// 注册路由
	router.GET("/1/g/1.0/example/:id/", func(c *gin.Context) {
		c.String(http.StatusOK, "不需要校验权限")
	})
	router.GET("/1/g/1.0/example/:id/edit/", func(c *gin.Context) {
		c.String(http.StatusOK, "userType校验通过")
	})
	router.GET("/1/g/1.0/example/:id/default/", func(c *gin.Context) {
		c.String(http.StatusOK, "redis权限校验通过")
	})

	t.Run("不需要校验权限", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/1/g/1.0/example/1/", nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, "不需要校验权限", w.Body.String())
	})
	t.Run("userType校验通过", func(t *testing.T) {
		auth := common.Auth{
			Exp:      jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)),
			UserType: "labor_cdl",
		}
		token, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
		assert.NoError(t, err)
		req, err := http.NewRequest("GET", "/1/g/1.0/example/1/edit/", nil)
		assert.NoError(t, err)
		req.Header.Set("HTTP_AUTHORIZATION", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, "userType校验通过", w.Body.String())
	})
	t.Run("userType校验不通过", func(t *testing.T) {
		auth := common.Auth{
			Exp:      jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)),
			UserType: "labor??", // 这里写了错误的类型
		}
		token, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
		assert.NoError(t, err)
		req, err := http.NewRequest("GET", "/1/g/1.0/example/1/edit/", nil)
		assert.NoError(t, err)
		req.Header.Set("HTTP_AUTHORIZATION", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.JSONEq(t, `{"detail": "您没有执行该操作的权限。"}`, w.Body.String())
	})
	t.Run("redis权限通过", func(t *testing.T) {
		auth := common.Auth{
			Exp:        jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)),
			InstanceId: 1,
			QsApp:      1,
			DrfUser:    "张三",
			Uid:        1,
			UserId:     1,
		}
		token, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
		assert.NoError(t, err)
		permissions := CachePermission{Actions: []string{"/1/g/1.0/example/:id/default/:GET"}}
		rbs, err := json.Marshal(permissions)
		assert.NoError(t, err)
		assert.NoError(t, ctx.Redis().Set(ctx.Context(), fmt.Sprintf("actions|Bearer %s", token), rbs, time.Duration(config.C.JWT.Exp)*time.Second).Err())
		t.Cleanup(func() {
			ctx.Redis().Del(ctx.Context(), fmt.Sprintf("actions|Bearer %s", token))
		})
		req, err := http.NewRequest("GET", "/1/g/1.0/example/1/default/", nil)
		assert.NoError(t, err)
		req.Header.Set("HTTP_AUTHORIZATION", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, "redis权限校验通过", w.Body.String())
	})
	t.Run("redis权限不通过", func(t *testing.T) {
		auth := common.Auth{
			Exp:        jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)),
			InstanceId: 1,
			QsApp:      1,
			DrfUser:    "张三",
			Uid:        1,
			UserId:     1,
		}
		token, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
		assert.NoError(t, err)
		permissions := CachePermission{Actions: []string{"/1/g/1.0/example/:id/default/:POST"}} // 这里更改了权限
		rbs, err := json.Marshal(permissions)
		assert.NoError(t, err)
		assert.NoError(t, ctx.Redis().Set(ctx.Context(), fmt.Sprintf("actions|%s", token), rbs, time.Duration(config.C.JWT.Exp)*time.Second).Err())
		t.Cleanup(func() {
			ctx.Redis().Del(ctx.Context(), fmt.Sprintf("actions|Bearer %s", token))
		})
		req, err := http.NewRequest("GET", "/1/g/1.0/example/1/default/", nil)
		assert.NoError(t, err)
		req.Header.Set("HTTP_AUTHORIZATION", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.JSONEq(t, `{"detail": "您没有执行该操作的权限。"}`, w.Body.String())
	})
}
