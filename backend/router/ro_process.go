package router

import (
	v1 "ConDetect/backend/app/api/v1"
	"ConDetect/backend/middleware"

	"github.com/gin-gonic/gin"
)

type ProcessRouter struct {
}

func (f *ProcessRouter) InitRouter(Router *gin.RouterGroup) {
	processRouter := Router.Group("process")
	processRouter.Use(middleware.JwtAuth()).Use(middleware.SessionAuth()).Use(middleware.PasswordExpired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		processRouter.GET("/ws", baseApi.ProcessWs)
		processRouter.POST("/stop", baseApi.StopProcess)
	}
}
