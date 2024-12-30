package migration

import (
	"ConDetect/backend/global"
	"ConDetect/backend/init/migration/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Init() {
	m := gormigrate.New(global.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.AddTableOperationLog,
		migrations.AddTableHost,
		migrations.AddTableMonitor,
		migrations.AddTableSetting,
		migrations.AddTableCronjob,
		migrations.AddTableImageRepo,
		migrations.AddDefaultGroup,
		migrations.UpdateTableHost,
		migrations.AddEntranceAndSSL,
		migrations.UpdateTableSetting,
		migrations.AddBindAndAllowIPs,
		migrations.UpdateCronjobWithSecond,
		migrations.AddMfaInterval,
		migrations.EncryptHostPassword,

		migrations.AddDefaultNetwork,
		migrations.UpdateTag,

		migrations.AddFavorite,
		migrations.AddBindAddress,
		migrations.AddCommandGroup,
		migrations.AddAppSyncStatus,

		migrations.AddDockerSockPath,
		migrations.AddSettingRecycleBin,

		migrations.AddXpackHideMenu,
		migrations.AddCronjobCommand,
		migrations.NewMonitorDB,
		migrations.AddNoAuthSetting,
		migrations.UpdateXpackHideMenu,
		migrations.AddMenuTabsSetting,
		migrations.AddDeveloperSetting,

		migrations.AddRedisCommand,
		migrations.AddMonitorMenu,
		migrations.AddProxy,
		migrations.AddCronJobColumn,
		migrations.AddForward,
		migrations.AddTrivy,
		migrations.AddSandfly,
		migrations.AddMeasureRecord,
		migrations.AddAlertMenu,
		migrations.AddComposeColumn,
	})
	if err := m.Migrate(); err != nil {
		global.LOG.Error(err)
		panic(err)
	}
	global.LOG.Info("Migration run successfully")
}
