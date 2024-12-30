package service

import (
	"ConDetect/backend/app/dto"
	"ConDetect/backend/app/model"
	"ConDetect/backend/constant"
	"ConDetect/backend/global"
	"ConDetect/backend/utils/common"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/aquasecurity/trivy/pkg/commands/artifact"
	"github.com/aquasecurity/trivy/pkg/flag"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

const (
	trivyResultDir = "trivy"
)

type VulnService struct{}

type IVulnService interface {
	Create(req dto.TrivyCreate) error
	Update(req dto.TrivyUpdate) error
	Delete(req dto.TrivyDelete) error
	HandleOnce(req dto.OperateByID) error
	CleanRecord(req dto.OperateByID) error
	SearchWithPage(req dto.SearchTrivyWithPage) (int64, interface{}, error)
	LoadTrivyRecords(req dto.TrivyRecordSearch) (int64, interface{}, error)
}

func NewIVulnService() IVulnService {
	return &VulnService{}
}
func (c *VulnService) CleanRecord(req dto.OperateByID) error {
	trivy, _ := vulnRepo.Get(commonRepo.WithByID(req.ID))
	if trivy.ID == 0 {
		return constant.ErrRecordNotFound
	}
	pathItem := path.Join(global.CONF.System.DataDir, trivyResultDir, trivy.Name)
	_ = os.RemoveAll(pathItem)
	return nil
}

func (c *VulnService) LoadTrivyRecords(req dto.TrivyRecordSearch) (int64, interface{}, error) {
	trivy, _ := vulnRepo.Get(commonRepo.WithByID(req.TrivyID))
	if trivy.ID == 0 {
		return 0, nil, constant.ErrRecordNotFound
	}
	jsonPaths := loadJsonByName(trivy.Name)
	if len(jsonPaths) == 0 {
		return 0, nil, nil
	}
	var filterFiles []string
	nyc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
	for _, item := range jsonPaths {
		t1, err := time.ParseInLocation(constant.DateTimeSlimLayout, item, nyc)
		if err != nil {
			continue
		}
		if t1.After(req.StartTime) && t1.Before(req.EndTime) {
			filterFiles = append(filterFiles, item)
		}
	}
	if len(filterFiles) == 0 {
		return 0, nil, nil
	}

	sort.Slice(filterFiles, func(i, j int) bool {
		return filterFiles[i] > filterFiles[j]
	})

	var records []string
	total, start, end := len(filterFiles), (req.Page-1)*req.PageSize, req.Page*req.PageSize
	if start > total {
		records = make([]string, 0)
	} else {
		if end >= total {
			end = total
		}
		records = filterFiles[start:end]
	}

	var datas []dto.TrivyRecord
	for i := 0; i < len(records); i++ {
		item := loadResultFromJson(path.Join(global.CONF.System.DataDir, trivyResultDir, trivy.Name, records[i]))
		datas = append(datas, item)
	}
	return int64(total), datas, nil
}
func loadResultFromJson(pathItem string) dto.TrivyRecord {
	var data dto.TrivyRecord
	data.Name = path.Base(pathItem)
	jsonFile, err := os.ReadFile(pathItem)
	if err != nil {
		data.Status = "Waiting"
		return data
	}
	var report types.Report
	if err := json.Unmarshal(jsonFile, &report); err != nil {
		data.Status = "Waiting"
		return data
	}
	if len(report.Results) == 0 {
		return data
	}
	data.Target = report.Results[0].Target
	data.CreatedAt = report.CreatedAt
	data.TotalCVE = 0
	for _, vuln := range report.Results[0].Vulnerabilities {
		var cve dto.Vulnerability
		cve.Title = vuln.Title
		cve.Description = vuln.Description
		cve.VulnerabilityID = vuln.VulnerabilityID
		cve.PkgName = vuln.PkgName
		cve.PrimaryURL = vuln.PrimaryURL
		cve.Severity = vuln.Severity
		cve.Status = vuln.Status.String()
		cve.InstalledVersion = vuln.InstalledVersion
		cve.FixedVersion = vuln.FixedVersion
		data.Vulnerabilities = append(data.Vulnerabilities, cve)
		data.TotalCVE += 1
	}
	data.Status = "Done"
	return data
}

func (v *VulnService) Create(req dto.TrivyCreate) error {
	trivy, _ := vulnRepo.Get(commonRepo.WithByName(req.Name))
	if trivy.ID != 0 {
		return constant.ErrRecordExist
	}
	if err := copier.Copy(&trivy, &req); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	if err := vulnRepo.Create(&trivy); err != nil {
		return err
	}
	return nil
}
func (v *VulnService) Update(req dto.TrivyUpdate) error {
	trivy, _ := vulnRepo.Get(commonRepo.WithByName(req.Name))
	if trivy.ID == 0 {
		return constant.ErrRecordNotFound
	}
	var trivyItem model.Trivy
	if err := copier.Copy(&trivyItem, &req); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	upMap := map[string]interface{}{}
	upMap["name"] = req.Name
	upMap["type"] = req.Type
	upMap["parallel"] = req.Parallel
	upMap["target"] = req.Target
	upMap["description"] = req.Description
	if err := vulnRepo.Update(req.ID, upMap); err != nil {
		return err
	}
	return nil
}
func (v *VulnService) Delete(req dto.TrivyDelete) error {
	for _, id := range req.Ids {
		trivy, _ := vulnRepo.Get(commonRepo.WithByID(id))
		if trivy.ID == 0 {
			continue
		}
		if req.RemoveRecord {
			_ = os.RemoveAll(path.Join(global.CONF.System.DataDir, trivyResultDir, trivy.Name))
		}
		if err := vulnRepo.Delete(commonRepo.WithByID(id)); err != nil {
			return err
		}
	}
	return nil
}
func (c *VulnService) SearchWithPage(req dto.SearchTrivyWithPage) (int64, interface{}, error) {
	total, commands, err := vulnRepo.Page(req.Page, req.PageSize, commonRepo.WithLikeName(req.Info), commonRepo.WithOrderRuleBy(req.OrderBy, req.Order))
	if err != nil {
		return 0, nil, err
	}
	var datas []dto.TrivyInfo
	for _, command := range commands {
		var item dto.TrivyInfo
		if err := copier.Copy(&item, &command); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		item.LastHandleDate = "-"
		datas = append(datas, item)
	}
	nyc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
	for i := 0; i < len(datas); i++ {
		logPaths := loadJsonByName(datas[i].Name)
		sort.Slice(logPaths, func(i, j int) bool {
			return logPaths[i] > logPaths[j]
		})
		if len(logPaths) != 0 {
			t1, err := time.ParseInLocation(constant.DateTimeSlimLayout, logPaths[0], nyc)
			if err != nil {
				continue
			}
			datas[i].LastHandleDate = t1.Format(constant.DateTimeLayout)
		}
	}

	return total, datas, err
}
func loadJsonByName(name string) []string {
	var logPaths []string
	pathItem := path.Join(global.CONF.System.DataDir, trivyResultDir, name)
	_ = filepath.Walk(pathItem, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || info.Name() == name {
			return nil
		}
		logPaths = append(logPaths, info.Name())
		return nil
	})
	return logPaths
}
func (v *VulnService) HandleOnce(req dto.OperateByID) error {
	trivy, _ := vulnRepo.Get(commonRepo.WithByID(req.ID))
	if trivy.ID == 0 {
		return constant.ErrRecordNotFound
	}
	timeNow := time.Now().Format(constant.DateTimeSlimLayout)
	report_path := path.Join(global.CONF.System.DataDir, trivyResultDir, trivy.Name, timeNow)
	if _, err := os.Stat(path.Dir(report_path)); err != nil {
		_ = os.MkdirAll(path.Dir(report_path), os.ModePerm)
	}
	go func() {
		switch trivy.Type {
		case 0:
			err := imageScan(trivy, report_path)
			if err != nil {
				global.LOG.Errorf("imagescan failed, err: %v", err)
			}
		case 1:
			err := fileSystemScan(trivy, report_path)
			if err != nil {
				global.LOG.Errorf("containerscan failed, err: %v", err)
			}
		case 2:
			err := imageScan(trivy, report_path)
			if err != nil {
				global.LOG.Errorf("filesystemscan failed, err: %v", err)
			}
		}
	}()
	return nil
}

func imageScan(trivy model.Trivy, report_path string) error {
	ctx := context.Background()
	DBFlagOptions, err := flag.NewDBFlagGroup().ToOptions()
	if err != nil {
		return fmt.Errorf("DBFlagOptions error: %v", err)
	}
	imageOptions, err := flag.NewImageFlagGroup().ToOptions()
	if err != nil {
		return fmt.Errorf("imageOptions error: %v", err)
	}
	target := []string{trivy.Target}
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
	cliOpt.ReportOptions.Output = report_path
	r, err := artifact.NewRunner(ctx, cliOpt)
	if err != nil {
		return fmt.Errorf("artifact.NewRunner err: %v", err)
	}
	defer r.Close(ctx)

	report, err := r.ScanImage(ctx, cliOpt)
	if err != nil {
		return fmt.Errorf("r.ScanImage err: %v", err)
	}
	err = r.Report(ctx, cliOpt, report)
	if err != nil {
		return fmt.Errorf("r.Report err: %v", err)
	}
	return nil
}

func fileSystemScan(trivy model.Trivy, report_path string) error {
	return nil
}
