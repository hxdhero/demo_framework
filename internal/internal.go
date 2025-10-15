package internal

import (
	"github.com/gin-gonic/gin"
	"lls_api/internal/handler"
	"lls_api/internal/middleware"
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/rerr"
)

func Servers() ([]app.Server, error) {
	// gin
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Recover())
	engine.Use(middleware.CORS())
	engine.Use(middleware.Log())
	engine.Use(middleware.DefaultPermission())

	// 第一个服务
	s1 := app.NewServer("lls_api", engine, config.C.Http.Host, config.C.Http.Port)
	// 基础handler
	base, err := handler.NewBaseHandler()
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	// base handler
	base.RegisterBaseService(s1)
	// 文件handler
	handler.NewBlueUsersaasHandler(base).RegisterHttpService(s1)
	handler.NewBlueUserHandler(base).RegisterHttpService(s1)
	handler.NewBluePlatformHandler(base).RegisterHttpService(s1)
	return []app.Server{*s1}, nil
}
