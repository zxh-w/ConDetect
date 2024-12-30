package router

import (
	v1 "ConDetect/backend/app/api/v1"
	"ConDetect/backend/middleware"

	"github.com/gin-gonic/gin"
)

type ToolboxRouter struct{}

func (s *ToolboxRouter) InitRouter(Router *gin.RouterGroup) {
	toolboxRouter := Router.Group("toolbox").
		Use(middleware.JwtAuth()).
		Use(middleware.SessionAuth()).
		Use(middleware.PasswordExpired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		toolboxRouter.POST("/device/base", baseApi.LoadDeviceBaseInfo)
		toolboxRouter.GET("/device/zone/options", baseApi.LoadTimeOption)
		toolboxRouter.POST("/device/update/conf", baseApi.UpdateDeviceConf)
		toolboxRouter.POST("/device/update/host", baseApi.UpdateDeviceHost)
		toolboxRouter.POST("/device/update/passwd", baseApi.UpdateDevicePasswd)
		toolboxRouter.POST("/device/update/swap", baseApi.UpdateDeviceSwap)
		toolboxRouter.POST("/device/update/byconf", baseApi.UpdateDeviceByFile)
		toolboxRouter.POST("/device/check/dns", baseApi.CheckDNS)
		toolboxRouter.POST("/device/conf", baseApi.LoadDeviceConf)

		toolboxRouter.POST("/scan", baseApi.ScanSystem)
		toolboxRouter.POST("/clean", baseApi.SystemClean)

		toolboxRouter.POST("/vuln/trivy/create", baseApi.CreateTrivy)
		toolboxRouter.POST("/vuln/trivy/update", baseApi.UpdateTrivy)
		toolboxRouter.POST("/vuln/trivy/search", baseApi.SearchTrivy)
		toolboxRouter.POST("/vuln/trivy/handle", baseApi.HandleTrivyScan)
		toolboxRouter.POST("/vuln/trivy/del", baseApi.DeleteTrivy)
		toolboxRouter.POST("/vuln/record/search", baseApi.SearchTrivyRecord)
		toolboxRouter.POST("/vuln/record/clean", baseApi.CleanTrivyRecord)

		toolboxRouter.POST("/baseline/check", baseApi.CheckBaseline)
		toolboxRouter.POST("/baseline/search", baseApi.SearchBaseline)

		toolboxRouter.POST("/macious/sandfly/create", baseApi.CreateSandfly)
		toolboxRouter.POST("/macious/sandfly/update", baseApi.UpdateSandfly)
		toolboxRouter.POST("/macious/sandfly/search", baseApi.SearchSandfly)
		toolboxRouter.POST("/macious/sandfly/handle", baseApi.HandleSandflyScan)
		toolboxRouter.POST("/macious/sandfly/del", baseApi.DeleteSandfly)
		toolboxRouter.POST("/macious/record/search", baseApi.SearchSandflyRecord)
		toolboxRouter.POST("/macious/record/clean", baseApi.CleanSandflyRecord)

		toolboxRouter.POST("/measure/container", baseApi.HandleContainerMeasure)
		toolboxRouter.POST("/measure/search", baseApi.SearchMeasureRecord)
	}
}
