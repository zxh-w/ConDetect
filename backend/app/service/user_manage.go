package service

import (
	"ConDetect/backend/app/dto"
	"ConDetect/backend/utils/cmd"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type UserManageService struct{}

type IUserManageService interface {
	// Create(req dto.) error
	Update(req dto.UserMmanage) error
	Delete(req dto.UserMmanage) error
	SearchWithName(req dto.UserSearch) (int64, interface{}, error)
}

func NewIUserManageService() IUserManageService {
	return &UserManageService{}
}

func (u *UserManageService) SearchWithName(req dto.UserSearch) (int64, interface{}, error) {

	file, err := os.Open("/etc/group")
	if err != nil {
		return 0, nil, fmt.Errorf("open group error: %v", err)
	}
	defer file.Close()
	docker_users := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if fields[0] == "docker" {
			for _, user := range strings.Split(fields[3], ",") {
				docker_users[user] = struct{}{}
			}
		}
	}
	file, err = os.Open("/etc/passwd")
	if err != nil {
		return 0, nil, fmt.Errorf("open pawsswd error: %v", err)
	}
	defer file.Close()
	var users []dto.UserMmanage
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) > 0 {
			var user dto.UserMmanage
			user.Name = fields[0]
			uid, _ := strconv.Atoi(fields[3])
			if uid < 4000 || uid > 65530 {
				continue
			}
			user.UID = fields[3]
			_, exit := docker_users[user.Name]
			if exit {
				user.Docker = true
			}
			s, _ := cmd.Exec(fmt.Sprintf("passwd -S %s", user.Name))
			if strings.Split(s, " ")[1] == "P" {
				user.Login = true
			}
			if strings.Contains(user.Name, req.Name) {
				users = append(users, user)
			}
		}
	}

	return 0, users, nil
}

func (u *UserManageService) Update(req dto.UserMmanage) error {
	file, err := os.Open("/etc/group")
	if err != nil {
		return fmt.Errorf("open group error: %v", err)
	}
	defer file.Close()
	docker_users := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if fields[0] == "docker" {
			for _, user := range strings.Split(fields[3], ",") {
				docker_users[user] = struct{}{}
			}
		}
	}
	_, exit := docker_users[req.Name]
	if exit != req.Docker {
		if req.Docker {
			_, err = cmd.Exec(fmt.Sprintf("usermod -aG docker %s", req.Name))
		} else {
			_, err = cmd.Exec(fmt.Sprintf("gpasswd -d %s docker", req.Name))
		}
		if err != nil {
			return fmt.Errorf("change user docker error: %v", err)
		}
	}
	s, _ := cmd.Exec(fmt.Sprintf("passwd -S %s", req.Name))
	login := strings.Split(s, " ")[1] == "P"
	if login != req.Login {
		if req.Login {
			_, err = cmd.Exec(fmt.Sprintf("usermod -U %v", req.Name))
		} else {
			_, err = cmd.Exec(fmt.Sprintf("usermod -L %v", req.Name))
		}
		if err != nil {
			return fmt.Errorf("change user login error: %v", err)
		}
	}
	return nil
}

func (u *UserManageService) Delete(req dto.UserMmanage) error {
	return nil
}
