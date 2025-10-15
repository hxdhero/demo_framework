package handler

import (
	"errors"
	"lls_api/common"
	"lls_api/internal/errs"
	"lls_api/internal/handler/dto"
	"lls_api/internal/model"
	"lls_api/pkg/app"
)

type BlueUsersaasHandler struct {
	*BaseHandler
}

func NewBlueUsersaasHandler(base *BaseHandler) *BlueUsersaasHandler {
	return &BlueUsersaasHandler{base}
}

func (h BlueUsersaasHandler) RegisterHttpService(s *app.Server) {
	v1 := s.Group("/1.0/usersaas")
	v1.PUT("/:id/status/", app.WrapHandler(h.Status))
}

func (h BlueUsersaasHandler) Status(ctx *app.Context) {
	id, err := ctx.Rest().GetUrlID()
	if err != nil {
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "请求参数错误", "2002509161818"))
		return
	}
	var req dto.ReqBlueUsersaasStatus
	if err := ctx.Rest().ShouldBind(&req); err != nil {
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "请求参数错误", "2002509161830"))
		return
	}
	if err := (&model.UserSaas{}).UpdateStatus(ctx, id, req.Status); err != nil {
		var disErr common.DisplayError
		if errors.As(err, &disErr) {
			ctx.Rest().HttpError(disErr)
			return
		}
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "修改usersaas状态错误", "2002509171614"))
		return
	}
	ctx.Rest().HttpSuccess("ok")
}
