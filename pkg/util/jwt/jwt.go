package jwt

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"lls_api/common"
	"lls_api/pkg/log"
	"lls_api/pkg/rerr"
	"time"
)

func NewJwtTime(t time.Time) jwt.NumericDate {
	return jwt.NumericDate{Time: t}
}

func ParseToken(tokenStr, secret string) (jwt.MapClaims, error) {
	// 解析 token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.DefaultContext().Infof("解析jwt token错误:%s", tokenStr)
			return nil, rerr.New("授权错误")
		}
		// 返回密钥
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, rerr.New("无效的授权")
	}
	return claims, nil
}

// GenerateToken 生成tokenStr
func GenerateToken(auth common.Auth, extra map[string]any, secret string) (string, error) {
	bs, err := json.Marshal(auth)
	if err != nil {
		return "", rerr.Wrap(err)
	}
	claims := make(jwt.MapClaims)
	if err := json.Unmarshal(bs, &claims); err != nil {
		return "", rerr.Wrap(err)
	}
	for k, v := range extra {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", rerr.Wrap(err)
	}
	return tokenStr, nil
}
