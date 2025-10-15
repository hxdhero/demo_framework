package handler

import (
	"fmt"
	"lls_api/common"
	"lls_api/internal/errs"
	"lls_api/internal/handler/dto"
	"lls_api/internal/model"
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/util/jwt"
	"time"
)

type BlueUserHandler struct {
	*BaseHandler
}

func NewBlueUserHandler(base *BaseHandler) *BlueUserHandler {
	return &BlueUserHandler{base}
}

func (h BlueUserHandler) RegisterHttpService(s *app.Server) {
	v1 := s.Group("/1.0/users")
	v1.POST("/login", app.WrapHandler(h.Login))
}

func (h BlueUserHandler) Login(ctx *app.Context) {
	var req dto.ReqBlueUserLogin
	if err := ctx.Rest().ShouldBind(&req); err != nil {
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "请求参数错误", "2002509191640"))
		return
	}
	iuser, err := (model.User{}).GetUser(ctx, req.Phone, req.QsApp, int32(req.RelatedId), "", 0)
	if err != nil {
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "获取用户失败", "2002509221008"))
		return
	}
	userType := iuser.UserType()
	instanceId := int(iuser.InstanceId())
	exp := time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)
	auth := common.Auth{
		Exp:        jwt.NewJwtTime(exp),
		InstanceId: instanceId,
		QsApp:      req.QsApp,
		DrfUser:    fmt.Sprintf("%s-%d", userType, instanceId),
		Uid:        int(iuser.GetID()),
		UserId:     instanceId,
		UserType:   userType,
		IsSuper:    iuser.IsSuper(),
		Extra:      nil,
	}
	// todo 这里可以封装一个方法 修改根据login instance来判断需要附加的内容
	tokenStr, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
	if err != nil {
		ctx.Log().Infof("生成token失败:%e", err)
		ctx.Rest().HttpError(errs.NewDisplayErr(ctx, err, "用户登录失败", "2002509221026"))
		return
	}
	resp := dto.RespBlueUserLogin{
		ExpAt:      auth.Exp,
		SaasID:     int(iuser.GetSaasID()),
		Tk:         tokenStr,
		Uid:        int(iuser.GetID()),
		UserSaasID: int(iuser.GetUserSaasID()),
	}
	// todo [用户登录][IP地址变更]
	// todo 记录首次登录信息和服务协议记录
	ctx.Rest().HttpSuccess(resp)
}
