package middleware

import "github.com/gin-gonic/gin"

func CORS() gin.HandlerFunc {
	// 定义允许的域名白名单
	allowOrigins := map[string]bool{
		"https://cdl-jie-h5.s2.iqusong.com": true,
		"https://jie.shzj178.com":           true,
		"http://worker.uat.shun178.com":     true,
		"https://worker.shzj178.com":        true,
		"http://saas.uat.shun178.com":       true,
		"https://saas.shzj178.com":          true,
	}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 如果请求来源在白名单中，则设置对应的 Access-Control-Allow-Origin
		if allowOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
