package pkg

import (
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/log"
	"lls_api/pkg/rdb"
	"lls_api/pkg/rerr"
	"lls_api/pkg/util/files"
	"lls_api/pkg/util/notification"
)

// InitFullDependence 初始化全部依赖
func InitFullDependence() {
	config.InitConfig()
	log.InitLog()
	app.InitDB()
	rdb.InitRedis()
	files.InitOss()
}

// InitGlobal 初始化全局依赖
func InitGlobal() []error {
	var initErrors []error
	// 初始化依赖
	InitFullDependence()

	if err := notification.InitLarkClient(); err != nil {
		initErrors = append(initErrors, rerr.WrapS(err, "初始化飞书失败"))
	}
	return initErrors
}
