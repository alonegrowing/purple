package middleware

import (
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	log "git.ur7.cn/opd/service/intersting/box/logging"
)

// 登陆校验 中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		// Process request
		c.Next()
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		handlerNamePathParts := strings.Split(c.HandlerName(), "/")
		handleName := handlerNamePathParts[len(handlerNamePathParts)-1]

		log.Infof("[ACCESS] %v | %20v | %3d | %13v | %15s | %-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			handleName,
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment)
	}
}
