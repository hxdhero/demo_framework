package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"lls_api/common"
	"lls_api/internal/errs"
	"lls_api/internal/model"
	"lls_api/internal/model/gen"
	"lls_api/pkg/app"
	"lls_api/pkg/log"
	"lls_api/pkg/utime"
	"net/http"
	"time"
)

// Auth 身份认证
func Auth() gin.HandlerFunc {
	return app.WrapHandler(func(ctx *app.Context) {
		// 处理用户信息
		if ctx.Session().AuthErr != nil {
			ctx.Rest().JSON(http.StatusUnauthorized, ctx.Session().AuthErr)
			ctx.Rest().Abort()
			return
		}
		t := ctx.Session().Auth.Exp
		if t.Before(time.Now()) {
			log.DefaultContext().Infof("token:%s 登录超时:%s", ctx.Session().AuthToken, t.Format(time.DateTime))
			ctx.Rest().JSON(http.StatusUnauthorized, errs.NewDisplayErr(ctx, nil, "登录超时，请重新登录", "2002509171744"))
			ctx.Rest().Abort()
			return
		}
		instance, err := model.LoginInstanceFromAuth(ctx, ctx.Session().Auth)
		if err != nil {
			log.DefaultContext().InfofSkip(1, "识别用户具体身份错误:%s", err.Error())
			ctx.Rest().JSON(http.StatusUnauthorized, errs.NewDisplayErr(ctx, err, "识别用户具体身份错误", "2002509301451"))
			ctx.Rest().Abort()
			return
		}

		if err := validateInstance(ctx, instance); err != nil {
			var disErr common.DisplayError
			if errors.As(err, &disErr) {
				ctx.Rest().JSON(http.StatusUnauthorized, disErr)
				ctx.Rest().Abort()
				return
			}
			ctx.Rest().JSON(http.StatusInternalServerError, err)
			ctx.Rest().Abort()
			return
		}

		ctx.Session().Instance = instance
	})
}

func validateInstance(ctx *app.Context, instance common.UserInstance) error {
	// 禁用员工不允许登录
	if labor, ok := instance.(model.Labor); ok {
		if !labor.IsEnable {
			return errs.NewDisplayErr(ctx, nil, "用户已禁用", "2002509301556")
		}
	}
	if qilianshe, ok := instance.(model.QiliansheUserClient); ok {
		var qlsClient model.QiliansheClient
		if err := ctx.DB().Take(&qlsClient, "id = ?", qilianshe.QlsClientID); err != nil {
			return err
		}
		if !qlsClient.QlsProbation.IsZero() && utime.DateBefore(qlsClient.QlsProbation, time.Now()) {
			ctx.Log().Infof("[用户禁用] 骑连社客户 7 天使用权已过:%+v", qilianshe.QlsClient)
			return errs.NewDisplayErr(ctx, nil, "用户已禁用", "2002509301557")
		}
	}
	if instance.GetStatus() == common.UserAgencyStatusDisabled {
		return errs.NewDisplayErr(ctx, nil, "用户已禁用", "2002509301558")
	}

	// 检查平台
	if userPlatform, ok := instance.(model.UserPlatform); ok {
		if !userPlatform.Platform.IsUse {
			return errs.NewDisplayErr(ctx, nil, "该集团已被禁用", "2002509301559")
		}
	}
	if userAgency, ok := instance.(model.UserAgency); ok {
		if !userAgency.Agency.Platform.IsUse {
			return errs.NewDisplayErr(ctx, nil, "该运营公司所属集团已被禁用", "2002509301601")
		}
	}
	if userCompany, ok := instance.(model.UserCompany); ok {
		var count int64
		if err := ctx.DB().Model(model.Company{}).
			Joins("Agency").
			Joins("Agency.Platform").
			Where(fmt.Sprintf("%s.id = ? and %s.is_use = ?", gen.TableNameBlueCompany, gen.TableNameBluePlatform), userCompany.CompanyID, true).
			Count(&count); err != nil {
			ctx.Log().Infof("数据库错误: %s", err.Error())
			return err
		}
		return errs.NewDisplayErr(ctx, nil, "该企业客户所属集团已被禁用", "2002509301602")
	}
	return nil
}
