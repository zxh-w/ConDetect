package service

import (
	"fmt"
	"os"
	"path"
	pathUtils "path"
	"strings"
	"time"

	"ConDetect/backend/app/model"
	"ConDetect/backend/constant"
	"ConDetect/backend/global"
	"ConDetect/backend/utils/cmd"
	"ConDetect/backend/utils/files"
	"ConDetect/backend/utils/ntp"

	"github.com/pkg/errors"
)

func (u *CronjobService) HandleJob(cronjob *model.Cronjob) {
	var (
		message []byte
		err     error
	)
	record := cronjobRepo.StartRecords(cronjob.ID, cronjob.KeepLocal, "")
	go func() {
		switch cronjob.Type {
		case "shell":
			if len(cronjob.Script) == 0 {
				return
			}
			record.Records = u.generateLogsPath(*cronjob, record.StartTime)
			_ = cronjobRepo.UpdateRecords(record.ID, map[string]interface{}{"records": record.Records})
			script := cronjob.Script
			if len(cronjob.ContainerName) != 0 {
				command := "sh"
				if len(cronjob.Command) != 0 {
					command = cronjob.Command
				}
				script = fmt.Sprintf("docker exec %s %s -c \"%s\"", cronjob.ContainerName, command, strings.ReplaceAll(cronjob.Script, "\"", "\\\""))
			}
			err = u.handleShell(cronjob.Type, cronjob.Name, script, record.Records)
			u.removeExpiredLog(*cronjob)
		case "curl":
			if len(cronjob.URL) == 0 {
				return
			}
			record.Records = u.generateLogsPath(*cronjob, record.StartTime)
			_ = cronjobRepo.UpdateRecords(record.ID, map[string]interface{}{"records": record.Records})
			err = u.handleShell(cronjob.Type, cronjob.Name, fmt.Sprintf("curl '%s'", cronjob.URL), record.Records)
			u.removeExpiredLog(*cronjob)
		case "ntp":
			err = u.handleNtpSync()
			u.removeExpiredLog(*cronjob)
		case "clean":
			messageItem := ""
			messageItem, err = u.handleSystemClean()
			message = []byte(messageItem)
			u.removeExpiredLog(*cronjob)
		}

		if err != nil {
			if len(message) != 0 {
				record.Records, _ = mkdirAndWriteFile(cronjob, record.StartTime, message)
			}
			cronjobRepo.EndRecords(record, constant.StatusFailed, err.Error(), record.Records)
			return
		}
		if len(message) != 0 {
			record.Records, err = mkdirAndWriteFile(cronjob, record.StartTime, message)
			if err != nil {
				global.LOG.Errorf("save file %s failed, err: %v", record.Records, err)
			}
		}
		cronjobRepo.EndRecords(record, constant.StatusSuccess, "", record.Records)
	}()
}

func (u *CronjobService) handleShell(cronType, cornName, script, logPath string) error {
	handleDir := fmt.Sprintf("%s/task/%s/%s", constant.DataDir, cronType, cornName)
	if _, err := os.Stat(handleDir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(handleDir, os.ModePerm); err != nil {
			return err
		}
	}
	if err := cmd.ExecCronjobWithTimeOut(script, handleDir, logPath, 24*time.Hour); err != nil {
		return err
	}
	return nil
}

func (u *CronjobService) handleNtpSync() error {
	ntpServer, err := settingRepo.Get(settingRepo.WithByKey("NtpSite"))
	if err != nil {
		return err
	}
	ntime, err := ntp.GetRemoteTime(ntpServer.Value)
	if err != nil {
		return err
	}
	if err := ntp.UpdateSystemTime(ntime.Format(constant.DateTimeLayout)); err != nil {
		return err
	}
	return nil
}

func handleTar(sourceDir, targetDir, name, exclusionRules string, secret string) error {
	if _, err := os.Stat(targetDir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return err
		}
	}

	excludes := strings.Split(exclusionRules, ",")
	excludeRules := ""
	for _, exclude := range excludes {
		if len(exclude) == 0 {
			continue
		}
		excludeRules += " --exclude " + exclude
	}
	path := ""
	if strings.Contains(sourceDir, "/") {
		itemDir := strings.ReplaceAll(sourceDir[strings.LastIndex(sourceDir, "/"):], "/", "")
		aheadDir := sourceDir[:strings.LastIndex(sourceDir, "/")]
		if len(aheadDir) == 0 {
			aheadDir = "/"
		}
		path += fmt.Sprintf("-C %s %s", aheadDir, itemDir)
	} else {
		path = sourceDir
	}

	commands := ""

	if len(secret) != 0 {
		extraCmd := "| openssl enc -aes-256-cbc -salt -k '" + secret + "' -out"
		commands = fmt.Sprintf("tar --warning=no-file-changed --ignore-failed-read --exclude-from=<(find %s -type s -print) -zcf %s %s %s %s", sourceDir, " -"+excludeRules, path, extraCmd, targetDir+"/"+name)
		global.LOG.Debug(strings.ReplaceAll(commands, fmt.Sprintf(" %s ", secret), "******"))
	} else {
		itemPrefix := pathUtils.Base(sourceDir)
		if itemPrefix == "/" {
			itemPrefix = ""
		}
		commands = fmt.Sprintf("tar --warning=no-file-changed --ignore-failed-read --exclude-from=<(find %s -type s -printf '%s' | sed 's|^|%s/|') -zcf %s %s %s", sourceDir, "%P\n", itemPrefix, targetDir+"/"+name, excludeRules, path)
		global.LOG.Debug(commands)
	}
	stdout, err := cmd.ExecWithTimeOut(commands, 24*time.Hour)
	if err != nil {
		if len(stdout) != 0 {
			global.LOG.Errorf("do handle tar failed, stdout: %s, err: %v", stdout, err)
			return fmt.Errorf("do handle tar failed, stdout: %s, err: %v", stdout, err)
		}
	}
	return nil
}

func handleUnTar(sourceFile, targetDir string, secret string) error {
	if _, err := os.Stat(targetDir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return err
		}
	}
	commands := ""
	if len(secret) != 0 {
		extraCmd := "openssl enc -d -aes-256-cbc -k '" + secret + "' -in " + sourceFile + " | "
		commands = fmt.Sprintf("%s tar -zxvf - -C %s", extraCmd, targetDir+" > /dev/null 2>&1")
		global.LOG.Debug(strings.ReplaceAll(commands, fmt.Sprintf(" %s ", secret), "******"))
	} else {
		commands = fmt.Sprintf("tar zxvfC %s %s", sourceFile, targetDir)
		global.LOG.Debug(commands)
	}

	stdout, err := cmd.ExecWithTimeOut(commands, 24*time.Hour)
	if err != nil {
		global.LOG.Errorf("do handle untar failed, stdout: %s, err: %v", stdout, err)
		return errors.New(stdout)
	}
	return nil
}

func backupLogFile(dstFilePath, websiteLogDir string, fileOp files.FileOp) error {
	if err := cmd.ExecCmd(fmt.Sprintf("tar -czf %s -C %s %s", dstFilePath, websiteLogDir, strings.Join([]string{"access.log", "error.log"}, " "))); err != nil {
		dstDir := path.Dir(dstFilePath)
		if err = fileOp.Copy(path.Join(websiteLogDir, "access.log"), dstDir); err != nil {
			return err
		}
		if err = fileOp.Copy(path.Join(websiteLogDir, "error.log"), dstDir); err != nil {
			return err
		}
		if err = cmd.ExecCmd(fmt.Sprintf("tar -czf %s -C %s %s", dstFilePath, dstDir, strings.Join([]string{"access.log", "error.log"}, " "))); err != nil {
			return err
		}
		_ = fileOp.DeleteFile(path.Join(dstDir, "access.log"))
		_ = fileOp.DeleteFile(path.Join(dstDir, "error.log"))
		return nil
	}
	return nil
}

func (u *CronjobService) handleSystemClean() (string, error) {
	return NewIDeviceService().CleanForCronjob()
}

func (u *CronjobService) removeExpiredLog(cronjob model.Cronjob) {
	records, _ := cronjobRepo.ListRecord(cronjobRepo.WithByJobID(int(cronjob.ID)), commonRepo.WithOrderBy("created_at desc"))
	if len(records) <= int(cronjob.RetainCopies) {
		return
	}
	for i := int(cronjob.RetainCopies); i < len(records); i++ {
		if len(records[i].File) != 0 {
			files := strings.Split(records[i].File, ",")
			for _, file := range files {
				_ = os.Remove(file)
			}
		}
		_ = cronjobRepo.DeleteRecord(commonRepo.WithByID(uint(records[i].ID)))
		_ = os.Remove(records[i].Records)
	}
}

func (u *CronjobService) generateLogsPath(cronjob model.Cronjob, startTime time.Time) string {
	dir := fmt.Sprintf("%s/task/%s/%s", constant.DataDir, cronjob.Type, cronjob.Name)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	path := fmt.Sprintf("%s/%s.log", dir, startTime.Format(constant.DateTimeSlimLayout))
	return path
}

func hasBackup(cronjobType string) bool {
	return cronjobType == "app" || cronjobType == "database" || cronjobType == "website" || cronjobType == "directory" || cronjobType == "snapshot" || cronjobType == "log"
}
