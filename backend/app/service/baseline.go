package service

import (
	"ConDetect/backend/app/dto"
	"ConDetect/backend/constant"
	"ConDetect/backend/utils/cmd"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type BaselineService struct{}

type IBaselineService interface {
	CheckBaseline() error
	SearchBaseline(dto.SearchBaselineWithPage) (int64, interface{}, error)
}

func NewIBaselineService() IBaselineService {
	return &BaselineService{}
}

func (v *BaselineService) CheckBaseline() error {

	script := path.Join(constant.DockerBenchDir, "docker-bench-security.sh -b -l docker-bench-security.log")

	err := cmd.ExecCmdWithDir(script, constant.DockerBenchDir)

	if err != nil {
		return fmt.Errorf("execScript DockerBench error: %v", err)
	}
	return nil
}

func (v *BaselineService) SearchBaseline(req dto.SearchBaselineWithPage) (int64, interface{}, error) {
	var baseline dto.BaselineInfo
	var dockerbench dto.DockerBench
	content, err := os.ReadFile(path.Join(constant.DockerBenchDir, "docker-bench-security.log.json"))
	if err != nil {
		return 0, baseline, fmt.Errorf("read docker-bench-security.log.json fail: %v", err)
	}
	err = json.Unmarshal(content, &dockerbench)
	if err != nil {
		return 0, baseline, fmt.Errorf("unmarshal docker-bench-security.log.json fail: %v", err)
	}
	baseline.Checks = dockerbench.Checks
	baseline.Score = dockerbench.Score
	baseline.Start = time.Unix(int64(dockerbench.Start), 0)
	baseline.End = time.Unix(int64(dockerbench.End), 0)

	var checkResults []dto.CheckResult
	var total int64
	for _, test := range dockerbench.Tests {
		var result dto.CheckResult
		for _, item := range test.Results {
			result.Result = item.Result
			if req.Level != "" && req.Level != result.Result {
				continue
			}
			result.Class = test.Desc
			result.Desc = item.Desc
			result.Id = item.Id
			result.Details = item.Details
			checkResults = append(checkResults, result)
			total += 1
		}
	}
	baseline.Results = checkResults
	startIndex := (req.Page - 1) * req.PageSize
	endIndex := startIndex + req.PageSize
	if endIndex > int(total) {
		endIndex = int(total)
	}
	return total, baseline.Results[startIndex:endIndex], nil
}
