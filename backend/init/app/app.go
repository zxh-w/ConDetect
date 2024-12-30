package app

import (
	"ConDetect/backend/utils/docker"
	"ConDetect/backend/utils/firewall"
	"path"

	"ConDetect/backend/constant"
	"ConDetect/backend/global"
	"ConDetect/backend/utils/files"
)

func Init() {
	constant.DataDir = global.CONF.System.DataDir
	constant.TrivyCacheDir = path.Join(constant.DataDir, "trivy-db")
	constant.DockerBenchDir = path.Join(constant.DataDir, "docker-bench-security")
	dirs := []string{constant.DataDir}

	fileOp := files.NewFileOp()
	for _, dir := range dirs {
		createDir(fileOp, dir)
	}

	_ = docker.CreateDefaultDockerNetwork()

	if f, err := firewall.NewFirewallClient(); err == nil {
		_ = f.EnableForward()
	}
}

func createDir(fileOp files.FileOp, dirPath string) {
	if !fileOp.Stat(dirPath) {
		_ = fileOp.CreateDir(dirPath, 0755)
	}
}
