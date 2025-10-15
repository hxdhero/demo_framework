package main

import (
	"context"
	"lls_api/internal"
	"lls_api/pkg"
	application "lls_api/pkg/app"
)

func main() {
	// 初始化应用依赖
	errs := pkg.InitGlobal()

	// 创建应用
	services, err := internal.Servers()
	if err != nil {
		panic(err)
	}
	app, err := application.NewApp(services)
	if err != nil {
		panic(err)
	}
	// 启动应用
	if err := app.Run(context.Background(), errs); err != nil {
		panic(err)
	}
}
