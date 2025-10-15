package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"lls_api/common"
	"lls_api/internal/model"
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/util/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newRequestWithAuth(auth common.Auth, method string, url string, body io.Reader) (*http.Request, error) {
	token, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
	if err != nil {
		return nil, err
	}
	return newRequestWithToken(token, method, url, body)
}

func newRequestWithToken(token string, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("HTTP_AUTHORIZATION", token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}

func TestCompanyZhaoshangSubAccounts(t *testing.T) {
	app.WithTestTransaction(t, ctx, func(ctx *app.Context) error {
		t.Run("user_company", func(t *testing.T) {
			// userCompany
			auth := common.Auth{
				Exp:        jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)),
				InstanceId: int(defaultUsersaas.InstanceId()),
				QsApp:      model.QS_APP_WEB_COMPANY,
				DrfUser:    "张三",
				Uid:        int(defaultUsersaas.GetID()),
				UserId:     int(defaultUsersaas.UserID),
			}
			req, err := newRequestWithAuth(auth, http.MethodGet, "/1/g/1.0/company_zhaoshang_sub_accounts/list_all", nil)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			service.ServeHTTP(w, req)
			assert.Equal(t, w.Code, 200, w.Body.String())
		})

		return nil
	})
}
