package router

import (
	v1 "ConDetect/backend/app/api/v1"
	"ConDetect/backend/middleware"

	"github.com/gin-gonic/gin"
)

type TerminalRouter struct{}

func (s *TerminalRouter) InitRouter(Router *gin.RouterGroup) {
	terminalRouter := Router.Group("terminals").
		Use(middleware.JwtAuth()).
		Use(middleware.SessionAuth()).
		Use(middleware.PasswordExpired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		terminalRouter.GET("", baseApi.WsSsh)
	}
}
