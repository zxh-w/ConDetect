package v1

import "ConDetect/backend/app/service"

type ApiGroup struct {
	BaseApi
}

var ApiGroupApp = new(ApiGroup)

var (
	authService      = service.NewIAuthService()
	dashboardService = service.NewIDashboardService()

	containerService       = service.NewIContainerService()
	composeTemplateService = service.NewIComposeTemplateService()
	imageRepoService       = service.NewIImageRepoService()
	imageService           = service.NewIImageService()
	dockerService          = service.NewIDockerService()

	hostService     = service.NewIHostService()
	groupService    = service.NewIGroupService()
	fileService     = service.NewIFileService()
	sshService      = service.NewISSHService()
	firewallService = service.NewIFirewallService()

	deviceService = service.NewIDeviceService()

	vulnService       = service.NewIVulnService()
	baselineService   = service.NewIBaselineService()
	maciousService    = service.NewIMaciousService()
	trustedServie     = service.NewITrustedService()
	userManageService = service.NewIUserManageService()

	settingService = service.NewISettingService()

	commandService = service.NewICommandService()

	logService     = service.NewILogService()
	processService = service.NewIProcessService()

	hostToolService = service.NewIHostToolService()

	recycleBinService = service.NewIRecycleBinService()
	favoriteService   = service.NewIFavoriteService()
)
