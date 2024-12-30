package constant

import (
	"path"

	"ConDetect/backend/global"
)

var (
	DataDir        = global.CONF.System.DataDir
	TrivyCacheDir  = path.Join(DataDir, "trivy-db")
	DockerBenchDir = path.Join(DataDir, "docker-bench-security")
	RecycleBinDir  = "/.condetect_clash"
)
