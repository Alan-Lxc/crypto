package middleware

import (
	"net/http"
	"regexp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域配置
func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie"}

	// 测试环境下模糊匹配本地开头的请求
	config.AllowOriginFunc = func(origin string) bool {
		if regexp.MustCompile(`^http://127\.0\.0\.1:\d+$`).MatchString(origin) {
			return true
		}
		if regexp.MustCompile(`^http://localhost:\d+$`).MatchString(origin) {
			return true
		}
		return false
	}

	config.AllowCredentials = true
	return cors.New(config)
}
func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		if context.Request.Method == http.MethodOptions {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
