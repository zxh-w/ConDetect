package hook

import (
	"encoding/base64"
	"strings"

	"ConDetect/backend/app/model"
	"ConDetect/backend/app/repo"
	"ConDetect/backend/constant"
	"ConDetect/backend/global"
	"ConDetect/backend/utils/cmd"
	"ConDetect/backend/utils/common"
	"ConDetect/backend/utils/encrypt"
)

func Init() {
	settingRepo := repo.NewISettingRepo()
	portSetting, err := settingRepo.Get(settingRepo.WithByKey("ServerPort"))
	if err != nil {
		global.LOG.Errorf("load service port from setting failed, err: %v", err)
	}
	global.CONF.System.Port = portSetting.Value
	ipv6Setting, err := settingRepo.Get(settingRepo.WithByKey("Ipv6"))
	if err != nil {
		global.LOG.Errorf("load ipv6 status from setting failed, err: %v", err)
	}
	global.CONF.System.Ipv6 = ipv6Setting.Value
	bindAddressSetting, err := settingRepo.Get(settingRepo.WithByKey("BindAddress"))
	if err != nil {
		global.LOG.Errorf("load bind address from setting failed, err: %v", err)
	}
	global.CONF.System.BindAddress = bindAddressSetting.Value
	sslSetting, err := settingRepo.Get(settingRepo.WithByKey("SSL"))
	if err != nil {
		global.LOG.Errorf("load service ssl from setting failed, err: %v", err)
	}
	global.CONF.System.SSL = sslSetting.Value

	OneDriveID, err := settingRepo.Get(settingRepo.WithByKey("OneDriveID"))
	if err != nil {
		global.LOG.Errorf("load onedrive info from setting failed, err: %v", err)
	}
	idItem, _ := base64.StdEncoding.DecodeString(OneDriveID.Value)
	global.CONF.System.OneDriveID = string(idItem)
	OneDriveSc, err := settingRepo.Get(settingRepo.WithByKey("OneDriveSc"))
	if err != nil {
		global.LOG.Errorf("load onedrive info from setting failed, err: %v", err)
	}
	scItem, _ := base64.StdEncoding.DecodeString(OneDriveSc.Value)
	global.CONF.System.OneDriveSc = string(scItem)

	if _, err := settingRepo.Get(settingRepo.WithByKey("SystemStatus")); err != nil {
		_ = settingRepo.Create("SystemStatus", "Free")
	}
	if err := settingRepo.Update("SystemStatus", "Free"); err != nil {
		global.LOG.Fatalf("init service before start failed, err: %v", err)
	}

	handleUserInfo(global.CONF.System.ChangeUserInfo, settingRepo)

	handleCronjobStatus()
}

func handleCronjobStatus() {
	_ = global.DB.Model(&model.JobRecords{}).Where("status = ?", constant.StatusWaiting).
		Updates(map[string]interface{}{
			"status":  constant.StatusFailed,
			"message": "the task was interrupted due to the restart of the condetect service",
		}).Error
}

func handleUserInfo(tags string, settingRepo repo.ISettingRepo) {
	if len(tags) == 0 {
		return
	}
	if tags == "all" {
		if err := settingRepo.Update("UserName", common.RandStrAndNum(10)); err != nil {
			global.LOG.Fatalf("init username before start failed, err: %v", err)
		}
		pass, _ := encrypt.StringEncrypt(common.RandStrAndNum(10))
		if err := settingRepo.Update("Password", pass); err != nil {
			global.LOG.Fatalf("init password before start failed, err: %v", err)
		}
		if err := settingRepo.Update("SecurityEntrance", common.RandStrAndNum(10)); err != nil {
			global.LOG.Fatalf("init entrance before start failed, err: %v", err)
		}
		return
	}
	if strings.Contains(global.CONF.System.ChangeUserInfo, "username") {
		if err := settingRepo.Update("UserName", common.RandStrAndNum(10)); err != nil {
			global.LOG.Fatalf("init username before start failed, err: %v", err)
		}
	}
	if strings.Contains(global.CONF.System.ChangeUserInfo, "password") {
		pass, _ := encrypt.StringEncrypt(common.RandStrAndNum(10))
		if err := settingRepo.Update("Password", pass); err != nil {
			global.LOG.Fatalf("init password before start failed, err: %v", err)
		}
	}
	if strings.Contains(global.CONF.System.ChangeUserInfo, "entrance") {
		if err := settingRepo.Update("SecurityEntrance", common.RandStrAndNum(10)); err != nil {
			global.LOG.Fatalf("init entrance before start failed, err: %v", err)
		}
	}

	sudo := cmd.SudoHandleCmd()
	_, _ = cmd.Execf("%s sed -i '/CHANGE_USER_INFO=%v/d' /usr/local/bin/1pctl", sudo, global.CONF.System.ChangeUserInfo)
}
