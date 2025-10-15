package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"lls_api/internal/handler/dto"
	"lls_api/internal/model"
	"lls_api/internal/model/gen"
	"lls_api/pkg/app"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPlatform(t *testing.T) {
	app.WithTestTransaction(t, ctx, func(ctx *app.Context) error {
		now := time.Now()
		platform := model.Platform{BluePlatform: gen.BluePlatform{CreateAt: now, ModifyAt: now, Name: "测试某平台"}}
		ctx.DB().Create(&platform)
		req, err := NewReqWithUser(t, defaultUsersaas, http.MethodGet, fmt.Sprintf("/1/g/1.0/platforms/%d/other_settings/", platform.ID), nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		service.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		var res dto.RespPlatformOtherSettings
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, res.Name, "测试某平台")
		return nil
	})
}
