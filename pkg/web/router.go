package web

import (
	"io/ioutil"
	"purple/pkg/web/handler"
	"purple/pkg/web/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

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
