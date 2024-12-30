package migrations

import (
	"encoding/json"
	"strings"
	"time"

	"ConDetect/backend/app/dto"
	"ConDetect/backend/app/model"
	"ConDetect/backend/constant"
	"ConDetect/backend/global"
	"ConDetect/backend/utils/common"
	"ConDetect/backend/utils/encrypt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableOperationLog = &gormigrate.Migration{
	ID: "20200809-add-table-operation-log",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&model.OperationLog{}, &model.LoginLog{})
	},
}

var AddTableHost = &gormigrate.Migration{
	ID: "20200818-add-table-host",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Host{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.Group{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.Command{}); err != nil {
			return err
		}
		group := model.Group{
			Name: "default", Type: "host", IsDefault: true,
		}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		host := model.Host{
			Name: "localhost", Addr: "127.0.0.1", User: "root", Port: 22, AuthMode: "password", GroupID: group.ID,
		}
		if err := tx.Create(&host).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddTableMonitor = &gormigrate.Migration{
	ID: "20200905-add-table-monitor",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&model.MonitorBase{}, &model.MonitorIO{}, &model.MonitorNetwork{})
	},
}

var AddTableSetting = &gormigrate.Migration{
	ID: "20200908-add-table-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Setting{}); err != nil {
			return err
		}
		encryptKey := common.RandStr(16)
		if err := tx.Create(&model.Setting{Key: "UserName", Value: global.CONF.System.Username}).Error; err != nil {
			return err
		}
		global.CONF.System.EncryptKey = encryptKey
		pass, _ := encrypt.StringEncrypt(global.CONF.System.Password)
		if err := tx.Create(&model.Setting{Key: "Password", Value: pass}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "Email", Value: ""}).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Setting{Key: "PanelName", Value: "ConDetect"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "Language", Value: "zh"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "Theme", Value: "light"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Setting{Key: "SessionTimeout", Value: "86400"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "LocalTime", Value: ""}).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Setting{Key: "ServerPort", Value: global.CONF.System.Port}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "SecurityEntrance", Value: global.CONF.System.Entrance}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "JWTSigningKey", Value: common.RandStr(16)}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "EncryptKey", Value: encryptKey}).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Setting{Key: "ExpirationTime", Value: time.Now().AddDate(0, 0, 10).Format(constant.DateTimeLayout)}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "ExpirationDays", Value: "0"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "ComplexityVerification", Value: "enable"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "MFAStatus", Value: "disable"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "MFASecret", Value: ""}).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Setting{Key: "MonitorStatus", Value: "enable"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "MonitorStoreDays", Value: "7"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Setting{Key: "MessageType", Value: "none"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "EmailVars", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "WeChatVars", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "DingVars", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "SystemVersion", Value: global.CONF.System.Version}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "SystemStatus", Value: "Free"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "AppStoreVersion", Value: ""}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddTableCronjob = &gormigrate.Migration{
	ID: "20200921-add-table-cronjob",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&model.Cronjob{}, &model.JobRecords{})
	},
}

var AddTableImageRepo = &gormigrate.Migration{
	ID: "20201009-add-table-imagerepo",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.ImageRepo{}, &model.ComposeTemplate{}, &model.Compose{}); err != nil {
			return err
		}
		item := &model.ImageRepo{
			Name:        "Docker Hub",
			Protocol:    "https",
			DownloadUrl: "docker.io",
			Status:      constant.StatusSuccess,
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddDefaultGroup = &gormigrate.Migration{
	ID: "2023022-change-default-group",
	Migrate: func(tx *gorm.DB) error {
		defaultGroup := &model.Group{
			Name:      "默认",
			IsDefault: true,
			Type:      "website",
		}
		if err := tx.Create(defaultGroup).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Group{}).Where("name = ? AND type = ?", "default", "host").Update("name", "默认").Error; err != nil {
			return err
		}
		return tx.Migrator().DropTable("website_groups")
	},
}

var UpdateTableHost = &gormigrate.Migration{
	ID: "20230410-update-table-host",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Host{}); err != nil {
			return err
		}
		return nil
	},
}

var AddEntranceAndSSL = &gormigrate.Migration{
	ID: "20230414-add-entrance-and-ssl",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Model(&model.Setting{}).
			Where("key = ? AND value = ?", "SecurityEntrance", "onepanel").
			Updates(map[string]interface{}{"value": ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "SSLType", Value: "self"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "SSLID", Value: "0"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "SSL", Value: "disable"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var UpdateTableSetting = &gormigrate.Migration{
	ID: "20200516-update-table-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "AppStoreLastModified", Value: "0"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddBindAndAllowIPs = &gormigrate.Migration{
	ID: "20230517-add-bind-and-allow",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "BindDomain", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "AllowIPs", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "TimeZone", Value: common.LoadTimeZoneByCmd()}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "NtpSite", Value: "pool.ntp.org"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "MonitorInterval", Value: "5"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var UpdateCronjobWithSecond = &gormigrate.Migration{
	ID: "20200524-update-table-cronjob",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Cronjob{}); err != nil {
			return err
		}
		var jobs []model.Cronjob
		if err := tx.Where("exclusion_rules != ?", "").Find(&jobs).Error; err != nil {
			return err
		}
		for _, job := range jobs {
			if strings.Contains(job.ExclusionRules, ";") {
				newRules := strings.ReplaceAll(job.ExclusionRules, ";", ",")
				if err := tx.Model(&model.Cronjob{}).Where("id = ?", job.ID).Update("exclusion_rules", newRules).Error; err != nil {
					return err
				}
			}
		}
		return nil
	},
}

var AddMfaInterval = &gormigrate.Migration{
	ID: "20230625-add-mfa-interval",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "MFAInterval", Value: "30"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "SystemIP", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "OneDriveID", Value: "MDEwOTM1YTktMWFhOS00ODU0LWExZGMtNmU0NWZlNjI4YzZi"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "OneDriveSc", Value: "akpuOFF+YkNXOU1OLWRzS1ZSRDdOcG1LT2ZRM0RLNmdvS1RkVWNGRA=="}).Error; err != nil {
			return err
		}
		return nil
	},
}

var EncryptHostPassword = &gormigrate.Migration{
	ID: "20230703-encrypt-host-password",
	Migrate: func(tx *gorm.DB) error {
		var hosts []model.Host
		if err := tx.Where("1 = 1").Find(&hosts).Error; err != nil {
			return err
		}

		var encryptSetting model.Setting
		if err := tx.Where("key = ?", "EncryptKey").Find(&encryptSetting).Error; err != nil {
			return err
		}
		global.CONF.System.EncryptKey = encryptSetting.Value

		for _, host := range hosts {
			if len(host.Password) != 0 {
				pass, err := encrypt.StringEncrypt(host.Password)
				if err != nil {
					return err
				}
				if err := tx.Model(&model.Host{}).Where("id = ?", host.ID).Update("password", pass).Error; err != nil {
					return err
				}
			}
			if len(host.PrivateKey) != 0 {
				key, err := encrypt.StringEncrypt(host.PrivateKey)
				if err != nil {
					return err
				}
				if err := tx.Model(&model.Host{}).Where("id = ?", host.ID).Update("private_key", key).Error; err != nil {
					return err
				}
			}
			if len(host.PassPhrase) != 0 {
				pass, err := encrypt.StringEncrypt(host.PassPhrase)
				if err != nil {
					return err
				}
				if err := tx.Model(&model.Host{}).Where("id = ?", host.ID).Update("pass_phrase", pass).Error; err != nil {
					return err
				}
			}
		}
		return nil
	},
}

var AddDefaultNetwork = &gormigrate.Migration{
	ID: "20230928-add-default-network",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "DefaultNetwork", Value: "all"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "LastCleanTime", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "LastCleanSize", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "LastCleanData", Value: ""}).Error; err != nil {
			return err
		}
		return nil
	},
}

var UpdateTag = &gormigrate.Migration{
	ID: "20231008-update-tag",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Tag{}); err != nil {
			return err
		}
		return nil
	},
}

var AddFavorite = &gormigrate.Migration{
	ID: "20231020-add-favorite",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Favorite{}); err != nil {
			return err
		}
		return nil
	},
}

var AddBindAddress = &gormigrate.Migration{
	ID: "20231024-add-bind-address",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "BindAddress", Value: "0.0.0.0"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "Ipv6", Value: "disable"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddCommandGroup = &gormigrate.Migration{
	ID: "20231030-add-command-group",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Command{}); err != nil {
			return err
		}
		defaultCommand := &model.Group{IsDefault: true, Name: "默认", Type: "command"}
		if err := tx.Create(defaultCommand).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Command{}).Where("1 = 1").Update("group_id", defaultCommand.ID).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddAppSyncStatus = &gormigrate.Migration{
	ID: "20231103-update-table-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "AppStoreSyncStatus", Value: "SyncSuccess"}).Error; err != nil {
			return err
		}
		return nil
	},
}
var AddDockerSockPath = &gormigrate.Migration{
	ID: "20231128-add-docker-sock-path",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "DockerSockPath", Value: "unix:///var/run/docker.sock"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddSettingRecycleBin = &gormigrate.Migration{
	ID: "20231129-add-setting-recycle-bin",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "FileRecycleBin", Value: "enable"}).Error; err != nil {
			return err
		}
		return nil
	},
}
var AddXpackHideMenu = &gormigrate.Migration{
	ID: "20240328-add-xpack-hide-menu",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "XpackHideMenu", Value: "{\"id\":\"1\",\"label\":\"/xpack\",\"isCheck\":true,\"title\":\"xpack.menu\",\"children\":[{\"id\":\"2\",\"title\":\"xpack.waf.name\",\"path\":\"/xpack/waf/dashboard\",\"label\":\"Dashboard\",\"isCheck\":true},{\"id\":\"3\",\"title\":\"xpack.tamper.tamper\",\"path\":\"/xpack/tamper\",\"label\":\"Tamper\",\"isCheck\":true},{\"id\":\"4\",\"title\":\"xpack.setting.setting\",\"path\":\"/xpack/setting\",\"label\":\"XSetting\",\"isCheck\":true}]}"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddCronjobCommand = &gormigrate.Migration{
	ID: "20240403-add-cronjob-command",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Cronjob{}); err != nil {
			return err
		}
		return nil
	},
}

var NewMonitorDB = &gormigrate.Migration{
	ID: "20240408-new-monitor-db",
	Migrate: func(tx *gorm.DB) error {
		var (
			bases    []model.MonitorBase
			ios      []model.MonitorIO
			networks []model.MonitorNetwork
		)
		_ = tx.Find(&bases).Error
		_ = tx.Find(&ios).Error
		_ = tx.Find(&networks).Error

		if err := global.MonitorDB.AutoMigrate(&model.MonitorBase{}, &model.MonitorIO{}, &model.MonitorNetwork{}); err != nil {
			return err
		}
		_ = global.MonitorDB.Exec("DELETE FROM monitor_bases").Error
		_ = global.MonitorDB.Exec("DELETE FROM monitor_ios").Error
		_ = global.MonitorDB.Exec("DELETE FROM monitor_networks").Error

		if len(bases) != 0 {
			for i := 0; i <= len(bases)/200; i++ {
				var itemData []model.MonitorBase
				if 200*(i+1) <= len(bases) {
					itemData = bases[200*i : 200*(i+1)]
				} else {
					itemData = bases[200*i:]
				}
				if len(itemData) != 0 {
					if err := global.MonitorDB.Create(&itemData).Error; err != nil {
						return err
					}
				}
			}
		}
		if len(ios) != 0 {
			for i := 0; i <= len(ios)/200; i++ {
				var itemData []model.MonitorIO
				if 200*(i+1) <= len(ios) {
					itemData = ios[200*i : 200*(i+1)]
				} else {
					itemData = ios[200*i:]
				}
				if len(itemData) != 0 {
					if err := global.MonitorDB.Create(&itemData).Error; err != nil {
						return err
					}
				}
			}
		}
		if len(networks) != 0 {
			for i := 0; i <= len(networks)/200; i++ {
				var itemData []model.MonitorNetwork
				if 200*(i+1) <= len(networks) {
					itemData = networks[200*i : 200*(i+1)]
				} else {
					itemData = networks[200*i:]
				}
				if len(itemData) != 0 {
					if err := global.MonitorDB.Create(&itemData).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	},
}

var AddNoAuthSetting = &gormigrate.Migration{
	ID: "20240328-add-no-auth-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "NoAuthSetting", Value: "200"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var UpdateXpackHideMenu = &gormigrate.Migration{
	ID: "20240411-update-xpack-hide-menu",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Model(&model.Setting{}).Where("key", "XpackHideMenu").Updates(map[string]interface{}{"value": "{\"id\":\"1\",\"label\":\"/xpack\",\"isCheck\":true,\"title\":\"xpack.menu\",\"children\":[{\"id\":\"2\",\"title\":\"xpack.waf.name\",\"path\":\"/xpack/waf/dashboard\",\"label\":\"Dashboard\",\"isCheck\":true},{\"id\":\"3\",\"title\":\"xpack.tamper.tamper\",\"path\":\"/xpack/tamper\",\"label\":\"Tamper\",\"isCheck\":true},{\"id\":\"4\",\"title\":\"xpack.gpu.gpu\",\"path\":\"/xpack/gpu\",\"label\":\"GPU\",\"isCheck\":true},{\"id\":\"5\",\"title\":\"xpack.setting.setting\",\"path\":\"/xpack/setting\",\"label\":\"XSetting\",\"isCheck\":true}]}"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddMenuTabsSetting = &gormigrate.Migration{
	ID: "20240415-add-menu-tabs-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "MenuTabs", Value: "disable"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddDeveloperSetting = &gormigrate.Migration{
	ID: "20240423-add-developer-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "DeveloperMode", Value: "disable"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddRedisCommand = &gormigrate.Migration{
	ID: "20240515-add-redis-command",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.RedisCommand{}); err != nil {
			return err
		}
		return nil
	},
}

var AddMonitorMenu = &gormigrate.Migration{
	ID: "20240517-update-xpack-hide-menu",
	Migrate: func(tx *gorm.DB) error {
		var (
			setting model.Setting
			menu    dto.XpackHideMenu
		)
		tx.Model(&model.Setting{}).Where("key", "XpackHideMenu").First(&setting)
		if err := json.Unmarshal([]byte(setting.Value), &menu); err != nil {
			return err
		}
		menu.Children = append(menu.Children, dto.XpackHideMenu{
			ID:      "6",
			Title:   "xpack.monitor.name",
			Path:    "/xpack/monitor/dashboard",
			Label:   "MonitorDashboard",
			IsCheck: true,
		})
		data, err := json.Marshal(menu)
		if err != nil {
			return err
		}
		return tx.Model(&model.Setting{}).Where("key", "XpackHideMenu").Updates(map[string]interface{}{"value": string(data)}).Error
	},
}

var AddProxy = &gormigrate.Migration{
	ID: "20240528-add-proxy",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&model.Setting{Key: "ProxyType", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "ProxyUrl", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "ProxyPort", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "ProxyUser", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "ProxyPasswd", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.Setting{Key: "ProxyPasswdKeep", Value: ""}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddForward = &gormigrate.Migration{
	ID: "202400611-add-forward",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Forward{}); err != nil {
			return err
		}
		return nil
	},
}

var AddCronJobColumn = &gormigrate.Migration{
	ID: "20240524-add-cronjob-command",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Cronjob{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTrivy = &gormigrate.Migration{
	ID: "20241017-add-trivy",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Trivy{}); err != nil {
			return err
		}
		return nil
	},
}
var AddSandfly = &gormigrate.Migration{
	ID: "20241026-add-sandfly",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Sandfly{}); err != nil {
			return err
		}
		return nil
	},
}
var AddMeasureRecord = &gormigrate.Migration{
	ID: "20241027-add-measure",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.MeasureRecord{}); err != nil {
			return err
		}
		return nil
	},
}
var AddAlertMenu = &gormigrate.Migration{
	ID: "20240706-update-xpack-hide-menu",
	Migrate: func(tx *gorm.DB) error {
		var (
			setting model.Setting
			menu    dto.XpackHideMenu
		)
		tx.Model(&model.Setting{}).Where("key", "XpackHideMenu").First(&setting)
		if err := json.Unmarshal([]byte(setting.Value), &menu); err != nil {
			return err
		}
		menu.Children = append(menu.Children, dto.XpackHideMenu{
			ID:      "7",
			Title:   "xpack.alert.alert",
			Path:    "/xpack/alert/dashboard",
			Label:   "XAlertDashboard",
			IsCheck: true,
		})
		data, err := json.Marshal(menu)
		if err != nil {
			return err
		}
		return tx.Model(&model.Setting{}).Where("key", "XpackHideMenu").Updates(map[string]interface{}{"value": string(data)}).Error
	},
}

var AddComposeColumn = &gormigrate.Migration{
	ID: "20240906-add-compose-command",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Compose{}); err != nil {
			return err
		}
		return nil
	},
}
