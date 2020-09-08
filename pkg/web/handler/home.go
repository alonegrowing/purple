package handler

import (
	"github.com/gin-gonic/gin"
)

type HomePageHandler struct {
}

func NewHomePageHandler() *HomePageHandler {
	return &HomePageHandler{}
}

func (m *HomePageHandler) Get(r *gin.Context) {
	var (
		code int64 = 0
	)
	r.JSON(200, gin.H{
		"code":    code,
		"message": "success",
		"data": map[string]interface{}{
			"name": "levin",
			"age":  28,
		},
	})
}
