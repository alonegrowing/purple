package web

import (
	"github.com/gin-gonic/gin"
	"purple/middleware"
	"purple/web/handler"

)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	route := gin.Default()
	route.Use(middleware.Logger())

	route.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	memberHandler := handler.NewHomePageHandler()

	// 路由定义
	route.GET("/api/member", middleware.Auth(), memberHandler.Get)

	return route
}