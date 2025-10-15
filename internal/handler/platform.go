package handler

import (
	"lls_api/internal/errs"
	"lls_api/internal/handler/dto"
	"lls_api/internal/model"
	"lls_api/pkg/app"
)

type PlatformHandler struct {
	*BaseHandler
}

func (h PlatformHandler) RegisterHttpService(s *app.Server) {
	v1 := s.Group("/1.0/platforms")
	v1.GET("/:id/other_settings/", app.WrapHandler(h.OtherSettings))

}

func NewBluePlatformHandler(baseHandler *BaseHandler) *PlatformHandler {
	return &PlatformHandler{BaseHandler: baseHandler}
}

func (h PlatformHandler) OtherSettings(ctx *app.Context) {
	var resp dto.RespPlatformOtherSettings
	id, err := ctx.Rest().GetUrlID()
	if err != nil {
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "请求参数错误", "2002509261629"))
		return
	}
	var platform model.Platform
	if err := ctx.DB().ByID(id, &platform); err != nil {
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "获取平台信息错误", "2002509261718"))
		return
	}
	resp.FromPlatForm(platform)
	ctx.Rest().HttpSuccess(resp)
}
