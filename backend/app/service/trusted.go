package service

import (
	"ConDetect/backend/app/dto"
	"ConDetect/backend/app/model"
	"ConDetect/backend/constant"
	"ConDetect/backend/utils/cmd"
	"ConDetect/backend/utils/docker"
	"ConDetect/backend/utils/entropyscan"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/aquasecurity/trivy/pkg/commands/artifact"
	"github.com/aquasecurity/trivy/pkg/flag"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/docker/docker/api/types/container"
)

type TrustedService struct{}

type ITrustedService interface {
	MeasureContainer(req dto.OperationWithName) error
	SearchWithPage(req dto.SearchMeasureWithPage) (int64, interface{}, error)
}

func NewITrustedService() ITrustedService {
	return &TrustedService{}
}
func (t *TrustedService) MeasureContainer(req dto.OperationWithName) error {

	image, upperDir, err := getContainerInfo(req.Name)
	if err != nil {
		return fmt.Errorf("get container image/upperDir fail: %v", err)
	}
	// 基线核查
	var record model.MeasureRecord
	record.Name = req.Name
	bashStr := fmt.Sprintf("docker-bench-security.sh -b -i %s -l %s", req.Name, "./log/docker-bench-security.log")
	script := path.Join(constant.DockerBenchDir, bashStr)
	err = cmd.ExecCmdWithDir(script, constant.DockerBenchDir)
	if err != nil {
		return fmt.Errorf("docker-bench exec fail: %v", err)
	}
	var dockerbench dto.DockerBench
	content, err := os.ReadFile(path.Join(constant.DockerBenchDir, "log", "docker-bench-security.log.json"))
	if err != nil {
		return fmt.Errorf("read docker-bench-security.log.json fail: %v", err)
	}
	err = json.Unmarshal(content, &dockerbench)
	if err != nil {
		return fmt.Errorf("unmarshal docker-bench-security.log.json fail: %v", err)
	}

	record.Baseline = float64(dockerbench.Score)

	// fmt.Println("========================基线核查======================", record.Baseline)
	// 漏洞探测
	ctx := context.Background()
	DBFlagOptions, err := flag.NewDBFlagGroup().ToOptions()
	if err != nil {
		return fmt.Errorf("DBFlagOptions error: %v", err)
	}
	imageOptions, err := flag.NewImageFlagGroup().ToOptions()
	if err != nil {
		return fmt.Errorf("imageOptions error: %v", err)
	}
	target := []string{image}
	scanOptions, err := flag.NewScanFlagGroup().ToOptions(target)
	if err != nil {
		return fmt.Errorf("scanOptions error: %v", err)
	}
	vulnerabilityOptions, err := flag.NewVulnerabilityFlagGroup().ToOptions()
	if err != nil {
		return fmt.Errorf("vulnerabilityOptions error: %v", err)
	}
	cliOpt := flag.Options{
		ImageOptions:         imageOptions,
		ScanOptions:          scanOptions,
		DBOptions:            DBFlagOptions,
		VulnerabilityOptions: vulnerabilityOptions,
	}
	cliOpt.CacheDir = constant.TrivyCacheDir
	cliOpt.SkipDBUpdate = true
	cliOpt.SkipJavaDBUpdate = true
	cliOpt.ScanOptions.OfflineScan = true
	cliOpt.ReportOptions.Format = types.FormatJSON
	r, err := artifact.NewRunner(ctx, cliOpt)
	if err != nil {
		return fmt.Errorf("artifact.NewRunner err: %v", err)
	}
	defer r.Close(ctx)

	report, err := r.ScanImage(ctx, cliOpt)
	if err != nil {
		return fmt.Errorf("r.ScanImage err: %v", err)
	}
	record.Vuln = 0
	for _, vuln := range report.Results[0].Vulnerabilities {

		switch vuln.Severity {
		case "LOW":
			record.Vuln += 1
		case "MEDIUM":
			record.Vuln += 2
		case "HIGH":
			record.Vuln += 3
		case "CRITICAL":
			record.Vuln += 4
		}
	}
	// fmt.Println("========================漏洞探测======================", record.Vuln)
	// 恶意文件扫描
	// 只检查elf文件,不检查正在运行的进程
	elfOnly, procOnly := true, false
	results, err := entropyscan.AnalyzeEntropy("", upperDir, 7.7, elfOnly, procOnly)
	if err != nil {
		return fmt.Errorf("AnalyzeEntropy error: %v", err)
	}
	for _, t := range results {
		record.Macious += float64(t.Entropy)
	}
	// fmt.Println("========================恶意文件扫描======================", record.Macious)
	record.Measure = record.Baseline + 0.1*record.Vuln + record.Macious

	meaureRecord, _ := trustedRepo.Get(commonRepo.WithByName(req.Name))
	if meaureRecord.ID != 0 {
		upMap := map[string]interface{}{}
		upMap["name"] = record.Name
		upMap["baseline"] = record.Baseline
		upMap["macious"] = record.Macious
		upMap["vuln"] = record.Vuln
		upMap["measure"] = record.Measure
		if err := trustedRepo.Update(meaureRecord.ID, upMap); err != nil {
			return err
		}
	} else {
		if err := trustedRepo.Create(&record); err != nil {
			return err
		}
	}
	return nil
}

func getContainerInfo(containerName string) (string, string, error) {
	client, err := docker.NewDockerClient()
	if err != nil {
		return "", "", err
	}
	defer client.Close()
	ctx := context.Background()
	container, err := client.ContainerInspect(ctx, containerName)
	return container.Image, container.GraphDriver.Data["UpperDir"], err
}

func (t *TrustedService) SearchWithPage(req dto.SearchMeasureWithPage) (int64, interface{}, error) {
	var datas []dto.MeasureResult
	containers, err := getContainerNameList()
	if err != nil {
		return 0, datas, err
	}
	total := 0
	for _, containerName := range containers {
		var data dto.MeasureResult
		data.Name = strings.TrimPrefix(containerName, "/")
		measureRecord, _ := trustedRepo.Get(commonRepo.WithByName(data.Name))
		if measureRecord.ID != 0 {
			data.Baseline = strconv.FormatFloat(measureRecord.Baseline, 'f', 2, 32)
			data.Macious = strconv.FormatFloat(measureRecord.Macious, 'f', 2, 64)
			data.Vuln = strconv.FormatFloat(measureRecord.Vuln, 'f', 2, 64)
			data.LastHandleDate = measureRecord.UpdatedAt.Format("2006-01-02 15:04:05")
			if measureRecord.Measure > 16 {
				data.Measure = "untrusted"
			} else {
				data.Measure = "trusted"
			}
		} else {
			data.LastHandleDate = "-"
			data.Baseline = "-"
			data.Macious = "-"
			data.Vuln = "-"
			data.Measure = "unmeasure"
		}
		total += 1
		datas = append(datas, data)
	}
	return int64(total), datas, nil
}

func getContainerNameList() ([]string, error) {
	client, err := docker.NewDockerClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	containers, err := client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	var datas []string
	for _, container := range containers {
		for _, name := range container.Names {
			if len(name) != 0 {
				datas = append(datas, strings.TrimPrefix(name, "/"))
			}
		}
	}
	return datas, nil
}
