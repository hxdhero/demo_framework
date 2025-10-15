package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"lls_api/common"
	"lls_api/internal/middleware"
	"lls_api/internal/model"
	"lls_api/internal/model/gen"
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/util/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateStatus(t *testing.T) {

	app.WithTestTransaction(t, ctx, func(ctx *app.Context) error {

		auth := common.Auth{
			Exp:        jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)),
			InstanceId: int(defaultUsersaas.InstanceId()),
			QsApp:      1,
			DrfUser:    "张三",
			Uid:        int(defaultUsersaas.GetID()),
			UserId:     int(defaultUsersaas.UserID),
		}
		token, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
		assert.NoError(t, err)

		// 没有执行权限
		permissions := middleware.CachePermission{Actions: []string{"/no permission"}}
		rbs, err := json.Marshal(permissions)
		assert.NoError(t, err)
		assert.NoError(t, ctx.Redis().Set(ctx.Context(), fmt.Sprintf("actions|%s", token), rbs, time.Duration(config.C.JWT.Exp)*time.Second).Err())
		body := map[string]any{"status": 1}
		bs, err := json.Marshal(body)
		assert.NoError(t, err)
		req, err := http.NewRequest("PUT", fmt.Sprintf("/1/g/1.0/usersaas/%d/status/", defaultUsersaas.ID), bytes.NewReader(bs))
		assert.NoError(t, err)
		req.Header.Set("HTTP_AUTHORIZATION", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		service.ServeHTTP(w, req)
		require.Equal(t, 403, w.Code, w.Body.String())
		require.Equal(t, "{\"detail\":\"您没有执行该操作的权限。\"}", w.Body.String())
		permissions = middleware.CachePermission{Actions: []string{"/1/g/1.0/usersaas/:id/status/:PUT"}}
		rbs, err = json.Marshal(permissions)
		assert.NoError(t, err)
		assert.NoError(t, ctx.Redis().Set(ctx.Context(), fmt.Sprintf("actions|Bearer %s", token), rbs, time.Duration(config.C.JWT.Exp)*time.Second).Err())

		// 请求数据
		defaultUsersaas.Status = 0
		assert.NoError(t, ctx.DB().Update(&defaultUsersaas, "status").Error)
		req, err = http.NewRequest("PUT", fmt.Sprintf("/1/g/1.0/usersaas/%d/status/", defaultUsersaas.ID), bytes.NewReader(bs))
		assert.NoError(t, err)
		req.Header.Set("HTTP_AUTHORIZATION", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w = httptest.NewRecorder()
		service.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200, w.Body.String())

		// 测试使用db更新
		defaultUsersaas.Name = "测试update"
		assert.NoError(t, ctx.DB().Update(&defaultUsersaas, "name").Error)
		var readUser model.UserSaas
		assert.NoError(t, ctx.DB().ByID(defaultUsersaas.ID, &readUser))
		assert.Equal(t, readUser.Name, defaultUsersaas.Name)

		// redis 保存权限
		permissions = middleware.CachePermission{Actions: []string{"/1/g/1.0/usersaas/:id/status/:PUT"}}
		rbs, err = json.Marshal(permissions)
		assert.NoError(t, err)
		assert.NoError(t, ctx.Redis().Set(ctx.Context(), fmt.Sprintf("actions|%s", token), rbs, time.Duration(config.C.JWT.Exp)*time.Second).Err())
		// 请求数据
		body = map[string]any{"status": 1}
		bs, err = json.Marshal(body)
		assert.NoError(t, err)
		req, err = http.NewRequest("PUT", fmt.Sprintf("/1/g/1.0/usersaas/%d/status/", defaultUsersaas.ID), bytes.NewReader(bs))
		assert.NoError(t, err)
		req.Header.Set("HTTP_AUTHORIZATION", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w = httptest.NewRecorder()
		service.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)

		// super用户不允许修改
		now := time.Now()
		user := model.User{BlueUser: gen.BlueUser{Mobile: "18989898990", CreateAt: now, ModifyAt: now}}
		if err := user.Create(ctx); err != nil {
			assert.NoError(t, err)
		}
		saas := model.Saas{BlueSaas: gen.BlueSaas{WalletID: 1, CreateAt: now, ModifyAt: now}}
		if err := saas.Create(ctx); err != nil {
			assert.NoError(t, err)
		}
		superUsersaas := model.UserSaas{BlueUsersaas: gen.BlueUsersaas{UserID: user.ID, SaasID: saas.ID, Name: "super用户", IsSuper: true, CreateAt: now, ModifyAt: now}}
		ctx.Session().IsDebug = true
		if err := superUsersaas.Create(ctx); err != nil {
			assert.NoError(t, err)
		}
		// redis 保存权限
		permissions = middleware.CachePermission{Actions: []string{"/1/g/1.0/usersaas/:id/status/:PUT"}}
		rbs, err = json.Marshal(permissions)
		assert.NoError(t, err)
		assert.NoError(t, ctx.Redis().Set(ctx.Context(), fmt.Sprintf("actions|%s", token), rbs, time.Duration(config.C.JWT.Exp)*time.Second).Err())
		body = map[string]any{"status": 1}
		bs, err = json.Marshal(body)
		assert.NoError(t, err)
		req, err = http.NewRequest("PUT", fmt.Sprintf("/1/g/1.0/usersaas/%d/status/", superUsersaas.ID), bytes.NewReader(bs))
		assert.NoError(t, err)
		setToken(req, token)
		w = httptest.NewRecorder()
		service.ServeHTTP(w, req)
		assert.Equal(t, ResDisplayErr(t, w).StandardDetail.Msg, "无法修改的用户")
		assert.Equal(t, ResDisplayErr(t, w).StandardDetail.Code, "2002509291341")
		return nil
	})

}
