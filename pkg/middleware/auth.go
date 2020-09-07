package middleware

import (
	"purple/pkg/macro"

	"github.com/gin-gonic/gin"
)

// 登陆校验 中间件
func Auth() gin.HandlerFunc {
	return func(r *gin.Context) {
		id := r.Query("id")
		if id == "2" {
			r.JSON(200, gin.H{
				"code":    macro.STATUS_AUTH_FAILED,
				"message": macro.ERR_MSG[macro.STATUS_AUTH_FAILED],
				"data":    []int{},
			})
			r.Abort()
		}
		r.Next()
	}
}
