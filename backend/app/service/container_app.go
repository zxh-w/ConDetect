package service

import (
	"ConDetect/backend/app/dto"
	"ConDetect/backend/utils/cmd"
	"fmt"
	"strings"
)

func (u *ContainerService) SearchAppWithContainer(req dto.AppSearch) (int64, interface{}, error) {
	s, err := cmd.Exec(fmt.Sprintf("appm -I -c %s -l", req.Container))
	if err != nil {
		return 0, nil, fmt.Errorf("exec app command error: %v", err)
	}
	features := strings.Split(s, "\n")

	var apps []dto.AppManage
	var index, name, enable, status string

	for _, info := range features {
		if strings.HasPrefix(info, "App index") {
			fmt.Println("App index", strings.Split(info, ":")[1])
			index = strings.Split(info, ":")[1]
		}
		if strings.HasPrefix(info, "App name") {
			fmt.Println("App name", strings.Split(info, ":")[1])
			name = strings.Split(info, ":")[1]
		}
		if strings.HasPrefix(info, "Service enable") {
			fmt.Println("App enable", strings.Split(info, ":")[1])
			enable = strings.Split(info, ":")[1]
		}
		if strings.HasPrefix(info, "Service status") {
			fmt.Println("App status", strings.Split(info, ":")[1])
			status = strings.Split(info, ":")[1]
		}
		if strings.HasPrefix(info, "--") {
			var app dto.AppManage
			app.Index = strings.TrimSpace(index)
			app.Name = strings.TrimSpace(name)
			app.Container = req.Container
			if strings.TrimSpace(enable) == "yes" {
				app.Enable = true
			} else {
				app.Enable = false
			}
			if strings.TrimSpace(status) == "running" {
				app.Status = true
			} else {
				app.Status = false
			}
			apps = append(apps, app)
		}
	}
	return 0, apps, nil
}

func (u *ContainerService) UpdateAppEnable(req dto.AppManage) error {
	if req.Enable {
		_, err := cmd.Exec(fmt.Sprintf("appm -e -c %s -n %s", req.Container, req.Name))
		if err != nil {
			return fmt.Errorf("exec enable app command error: %v", err)
		}
	} else {
		_, err := cmd.Exec(fmt.Sprintf("appm -d -c %s -n %s", req.Container, req.Name))
		if err != nil {
			return fmt.Errorf("exec disenable app command error: %v", err)
		}
	}
	return nil
}

func (u *ContainerService) UpdateAppStatus(req dto.AppManage) error {
	if req.Status {
		_, err := cmd.Exec(fmt.Sprintf("appm -s -c %s -n %s", req.Container, req.Name))
		if err != nil {
			return fmt.Errorf("exec start app command error: %v", err)
		}
	} else {
		_, err := cmd.Exec(fmt.Sprintf("appm -S -c %s -n %s", req.Container, req.Name))
		if err != nil {
			return fmt.Errorf("exec stop app command error: %v", err)
		}
	}
	return nil
}

func (u *ContainerService) UninstallApp(req dto.AppManage) error {
	_, err := cmd.Exec(fmt.Sprintf("appm -u -c %s -n %s", req.Container, req.Name))
	if err != nil {
		return fmt.Errorf("exec uninstall app command error: %v", err)
	}
	return nil
}
