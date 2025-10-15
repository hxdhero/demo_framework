package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"lls_api/common"
	"lls_api/pkg"
	"lls_api/pkg/config"
	"lls_api/pkg/util/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	pkg.InitGlobal()

	router := gin.New()
	router.Use(Auth())
	router.GET("/test_auth/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	// token 过期
	auth := common.Auth{
		Exp: jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second * -1)),
	}
	tokenStr, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
	assert.NoError(t, err)
	// 创建请求
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	var displayErr common.DisplayError
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &displayErr))
	assert.Equal(t, displayErr.StandardDetail.Code, "2002509171741")
	assert.Equal(t, displayErr.StandardDetail.Msg, "授权已过期")
}
