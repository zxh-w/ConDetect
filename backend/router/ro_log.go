package router

import (
	v1 "ConDetect/backend/app/api/v1"
	"ConDetect/backend/middleware"

	"github.com/gin-gonic/gin"
)

type LogRouter struct{}

func (s *LogRouter) InitRouter(Router *gin.RouterGroup) {
	operationRouter := Router.Group("logs")
	operationRouter.Use(middleware.JwtAuth()).Use(middleware.SessionAuth()).Use(middleware.PasswordExpired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		operationRouter.POST("/login", baseApi.GetLoginLogs)
		operationRouter.POST("/operation", baseApi.GetOperationLogs)
		operationRouter.POST("/clean", baseApi.CleanLogs)
		operationRouter.GET("/system/files", baseApi.GetSystemFiles)
		operationRouter.POST("/system", baseApi.GetSystemLogs)
	}
}
