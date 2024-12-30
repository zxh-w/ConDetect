package service

import "ConDetect/backend/app/repo"

var (
	commonRepo = repo.NewCommonRepo()

	imageRepoRepo = repo.NewIImageRepoRepo()
	composeRepo   = repo.NewIComposeTemplateRepo()

	cronjobRepo = repo.NewICronjobRepo()

	hostRepo    = repo.NewIHostRepo()
	groupRepo   = repo.NewIGroupRepo()
	commandRepo = repo.NewICommandRepo()

	vulnRepo    = repo.NewIVulnRepo()
	maciousRepo = repo.NewIMaciousRepo()
	trustedRepo = repo.NewITrustedRepo()

	settingRepo = repo.NewISettingRepo()

	logRepo = repo.NewILogRepo()

	favoriteRepo = repo.NewIFavoriteRepo()
)
