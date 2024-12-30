package service

import (
	"ConDetect/backend/app/dto"
	"ConDetect/backend/app/model"
	"ConDetect/backend/buserr"
	"ConDetect/backend/constant"
	"ConDetect/backend/global"
	"ConDetect/backend/utils/cmd"
	"ConDetect/backend/utils/common"
	"ConDetect/backend/utils/entropyscan"
	"ConDetect/backend/utils/files"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

const (
	sandflyResultDir = "sandfly"
)

type MaciousService struct{}

type IMaciousService interface {
	Create(req dto.SandflyCreate) error
	Update(req dto.SandflyUpdate) error
	Delete(req dto.SandflyDelete) error
	HandleOnce(req dto.OperateByID) error
	CleanRecord(req dto.OperateByID) error
	SearchWithPage(req dto.SearchSandflyWithPage) (int64, interface{}, error)
	LoadSandflyRecords(req dto.SandflyRecordSearch) (int64, interface{}, error)
}

func NewIMaciousService() IMaciousService {
	return &MaciousService{}
}
func (c *MaciousService) CleanRecord(req dto.OperateByID) error {
	sandfly, _ := maciousRepo.Get(commonRepo.WithByID(req.ID))
	if sandfly.ID == 0 {
		return constant.ErrRecordNotFound
	}
	pathItem := path.Join(global.CONF.System.DataDir, sandflyResultDir, sandfly.Name)
	_ = os.RemoveAll(pathItem)
	return nil
}

func (c *MaciousService) LoadSandflyRecords(req dto.SandflyRecordSearch) (int64, interface{}, error) {
	sandfly, _ := maciousRepo.Get(commonRepo.WithByID(req.SandflyID))
	if sandfly.ID == 0 {
		return 0, nil, constant.ErrRecordNotFound
	}
	jsonPaths := loadJsonBySandfly(sandfly.Name)
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
	var datas []dto.SandflyRecord
	for i := 0; i < len(records); i++ {
		item := loadRecordFromJson(path.Join(global.CONF.System.DataDir, sandflyResultDir, sandfly.Name, records[i]))
		datas = append(datas, item)
	}
	return int64(total), datas, nil
}
func loadJsonBySandfly(name string) []string {
	var logPaths []string
	pathItem := path.Join(global.CONF.System.DataDir, sandflyResultDir, name)
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
func loadRecordFromJson(pathItem string) dto.SandflyRecord {
	var data dto.SandflyRecord
	data.Name = path.Base(pathItem)
	jsonFile, err := os.ReadFile(pathItem)
	if err != nil {
		data.Status = "Waiting"
		return data
	}
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		data.Status = "Waiting"
		return data
	}
	data.Status = "Done"
	return data
}
func (c *MaciousService) HandleOnce(req dto.OperateByID) error {
	sandfly, _ := maciousRepo.Get(commonRepo.WithByID(req.ID))
	if sandfly.ID == 0 {
		return constant.ErrRecordNotFound
	}
	if cmd.CheckIllegal(sandfly.Path) {
		return buserr.New(constant.ErrCmdIllegal)
	}
	var record dto.SandflyRecord
	record.Path = sandfly.Path
	record.CreatedAt = time.Now()
	timeNow := time.Now().Format(constant.DateTimeSlimLayout)
	jsonFile := path.Join(global.CONF.System.DataDir, sandflyResultDir, sandfly.Name, timeNow)
	if _, err := os.Stat(path.Dir(jsonFile)); err != nil {
		_ = os.MkdirAll(path.Dir(jsonFile), os.ModePerm)
	}
	record.Name = timeNow
	record.Status = "Done"
	go func() {
		// 只检查elf文件,不检查正在运行的进程
		elfOnly, procOnly := true, false
		results, err := entropyscan.AnalyzeEntropy("", sandfly.Path, 7.7, elfOnly, procOnly)
		if err != nil {
			global.LOG.Errorf("macious Sandfly scan failed, err: %v", err)
		}
		fs := files.NewFileOp()
		totalMacious := 0
		for _, result := range results {
			filePath := path.Join(result.Path, result.Name)
			switch sandfly.InfectedStrategy {
			case "remove":
				fs.DeleteFile(filePath)
			case "move":
				dir := path.Join(sandfly.InfectedDir, "condetect-infected", sandfly.Name, timeNow)
				if _, err := os.Stat(dir); err != nil {
					_ = os.MkdirAll(dir, os.ModePerm)
				}
				fs.Mv(filePath, path.Join(dir, result.Name))
			case "copy":
				dir := path.Join(sandfly.InfectedDir, "condetect-infected", sandfly.Name, timeNow)
				if _, err := os.Stat(dir); err != nil {
					_ = os.MkdirAll(dir, os.ModePerm)
				}
				fs.Copy(filePath, path.Join(dir, result.Name))
			}
			totalMacious += 1
		}
		record.Macious = results
		record.TotalMacious = uint(totalMacious)
		if err = writeSandflyRecord(record, jsonFile); err != nil {
			global.LOG.Errorf("macious Sandfly write failed, err: %v", err)
		}
	}()
	return nil
}

// Wtite results
func writeSandflyRecord(record dto.SandflyRecord, jsonFile string) error {
	// 创建或打开指定的 JSON 文件
	file, err := os.OpenFile(jsonFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create JSON file: %v", err)
	}
	defer file.Close()

	// 编码 FileData 结构体为 JSON 格式
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 设置缩进以美化输出
	if err := encoder.Encode(record); err != nil {
		return fmt.Errorf("failed to encode file data to JSON: %v", err)
	}
	return nil
}
func (c *MaciousService) SearchWithPage(req dto.SearchSandflyWithPage) (int64, interface{}, error) {
	total, commands, err := maciousRepo.Page(req.Page, req.PageSize, commonRepo.WithLikeName(req.Info), commonRepo.WithOrderRuleBy(req.OrderBy, req.Order))
	if err != nil {
		return 0, nil, err
	}
	var datas []dto.SandflyInfo
	for _, command := range commands {
		var item dto.SandflyInfo
		if err := copier.Copy(&item, &command); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		item.LastHandleDate = "-"
		datas = append(datas, item)
	}
	nyc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
	for i := 0; i < len(datas); i++ {
		logPaths := loadRecordByName(datas[i].Name)
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
func loadRecordByName(name string) []string {
	var logPaths []string
	pathItem := path.Join(global.CONF.System.DataDir, sandflyResultDir, name)
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
func (c *MaciousService) Create(req dto.SandflyCreate) error {
	sandfly, _ := maciousRepo.Get(commonRepo.WithByName(req.Name))
	if sandfly.ID != 0 {
		return constant.ErrRecordExist
	}
	if err := copier.Copy(&sandfly, &req); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	if sandfly.InfectedStrategy == "none" || sandfly.InfectedStrategy == "remove" {
		sandfly.InfectedDir = ""
	}
	if err := maciousRepo.Create(&sandfly); err != nil {
		return err
	}
	return nil
}

func (c *MaciousService) Update(req dto.SandflyUpdate) error {
	sandfly, _ := maciousRepo.Get(commonRepo.WithByName(req.Name))
	if sandfly.ID == 0 {
		return constant.ErrRecordNotFound
	}
	if req.InfectedStrategy == "none" || req.InfectedStrategy == "remove" {
		req.InfectedDir = ""
	}
	var sandflyItem model.Sandfly
	if err := copier.Copy(&sandflyItem, &req); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	upMap := map[string]interface{}{}
	upMap["name"] = req.Name
	upMap["path"] = req.Path
	upMap["infected_dir"] = req.InfectedDir
	upMap["infected_strategy"] = req.InfectedStrategy
	upMap["description"] = req.Description
	if err := maciousRepo.Update(req.ID, upMap); err != nil {
		return err
	}
	return nil
}
func (c *MaciousService) Delete(req dto.SandflyDelete) error {
	for _, id := range req.Ids {
		sandfly, _ := maciousRepo.Get(commonRepo.WithByID(id))
		if sandfly.ID == 0 {
			continue
		}
		if req.RemoveRecord {
			_ = os.RemoveAll(path.Join(global.CONF.System.DataDir, sandflyResultDir, sandfly.Name))
		}
		if req.RemoveInfected {
			_ = os.RemoveAll(path.Join(sandfly.InfectedDir, "condetect-infected", sandfly.Name))
		}
		if err := maciousRepo.Delete(commonRepo.WithByID(id)); err != nil {
			return err
		}
	}
	return nil
}
