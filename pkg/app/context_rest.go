package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	jjwt "github.com/golang-jwt/jwt/v5"
	"lls_api/common"
	"lls_api/pkg/config"
	"lls_api/pkg/log"
	"lls_api/pkg/util/jwt"
	"lls_api/pkg/util/uuid"
	"strings"
)

/*
rest服务相关
*/

func ContextWithGin(c *gin.Context) *Context {
	requestID := c.GetHeader("HTTP_X_REQUEST_ID")
	if requestID == "" {
		requestID = uuid.UUID4()
	}
	// 处理requestID
	requestID = "#" + strings.Split(requestID, "#")[len(strings.Split(requestID, "#"))-1]

	// / 处理token
	tokenStr := extractToken(c)
	auth, authErr := getAuth(tokenStr, requestID, config.C.JWT.Secret)
	loggerContext := &log.LoggerContext{
		RequestID: requestID,
		User:      "unKnow",
	}
	if authErr == nil && auth.DrfUser != "" {
		loggerContext.User = auth.DrfUser
	}

	// 会话信息
	sessionInfo := &SessionInfo{
		RequestID: requestID,
		Auth:      auth,
		AuthErr:   authErr,
		AuthToken: tokenStr,
	}

	appContext := NewAppContext(context.Background(), c, loggerContext, sessionInfo)

	return appContext
}

func getAuth(tokenStr string, requestID string, secret string) (common.Auth, error) {
	var auth common.Auth
	claim, err := jwt.ParseToken(tokenStr, secret)
	if err != nil {
		if errors.Is(err, jjwt.ErrTokenExpired) {
			log.DefaultContext().Infof("token:%s 授权已过期", tokenStr)
			return auth, common.DisplayError{
				OriginErr: err,
				StandardDetail: common.StandardDetail{
					Msg:        "授权已过期",
					Code:       "2002509171741",
					XRequestID: requestID,
				},
			}
		}

		log.DefaultContext().Infof("token:%s 解析用户身份错误", tokenStr)
		return auth, common.DisplayError{
			OriginErr: err,
			StandardDetail: common.StandardDetail{
				Msg:        "解析用户身份错误",
				Code:       "2002509171739",
				XRequestID: requestID,
			},
		}
	}

	bs, err := json.Marshal(claim)
	if err != nil {
		log.DefaultContext().Infof("token:%s marshal token err:%e", tokenStr, err)
		return auth, common.DisplayError{
			OriginErr: err,
			StandardDetail: common.StandardDetail{
				Msg:        "授权格式错误",
				Code:       "2002509171747",
				XRequestID: requestID,
			},
		}
	}

	if err := json.Unmarshal(bs, &auth); err != nil {
		log.DefaultContext().Infof("tokenStr:%s unmarshal authuser err:%e", tokenStr, err)
		return auth, common.DisplayError{
			OriginErr: err,
			StandardDetail: common.StandardDetail{
				Msg:        "授权格式错误",
				Code:       "2002509171747",
				XRequestID: requestID,
			},
		}
	}
	return auth, nil
}

func extractToken(c *gin.Context) string {
	// 从 Authorization Header 获取
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		authHeader = c.GetHeader("HTTP_AUTHORIZATION")
	}
	if authHeader != "" {
		// 支持 Bearer 和直接 token 两种格式
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimSpace(authHeader[7:])
		}
		return strings.TrimSpace(authHeader)
	}

	// 从 query 参数获取
	if token := c.Query("token"); token != "" {
		return token
	}

	// 从表单获取
	if token := c.PostForm("token"); token != "" {
		return token
	}

	// 从 cookie 获取
	if token, err := c.Cookie("token"); err == nil && token != "" {
		return token
	}

	return ""
}
