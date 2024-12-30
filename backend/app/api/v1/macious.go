package v1

import (
	"ConDetect/backend/app/api/v1/helper"
	"ConDetect/backend/app/dto"
	"ConDetect/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags Macious
// @Summary Handle sandfly scan
// @Description 执行恶意文件扫描扫描
// @Accept json
// @Param request body dto.OperateByID true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/macious/sandfly/handle [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":true,"db":"sandflys","output_column":"name","output_value":"name"}],"formatZH":"执行恶意文件扫描 [name]","formatEN":"handle sandfly scan [name]"}
func (b *BaseApi) HandleSandflyScan(c *gin.Context) {
	var req dto.OperateByID
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := maciousService.HandleOnce(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Macious
// @Summary Create sandfly
// @Description 创建扫描规则
// @Accept json
// @Param request body dto.SandflyCreate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/macious/sandfly/create [post]
// @x-panel-log {"bodyKeys":["name","path"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建扫描规则 [name][path]","formatEN":"create sandfly [name][path]"}
func (b *BaseApi) CreateSandfly(c *gin.Context) {
	var req dto.SandflyCreate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := maciousService.Create(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Macious
// @Summary Update sandfly
// @Description 修改扫描规则
// @Accept json
// @Param request body dto.SandflyUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/macious/sandfly/update [post]
// @x-panel-log {"bodyKeys":["name","path"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"修改恶意文件扫描规则 [name][path]","formatEN":"update sandfly [name][path]"}
func (b *BaseApi) UpdateSandfly(c *gin.Context) {
	var req dto.SandflyUpdate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := maciousService.Update(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Macious
// @Summary Delete sandfly
// @Description 删除扫描规则
// @Accept json
// @Param request body dto.SandflyDelete true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/macious/sandfly/del [post]
// @x-panel-log {"bodyKeys":["ids"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"ids","isList":true,"db":"sandflys","output_column":"name","output_value":"names"}],"formatZH":"删除恶意文件扫描规则 [names]","formatEN":"delete trivy [names]"}
func (b *BaseApi) DeleteSandfly(c *gin.Context) {
	var req dto.SandflyDelete
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := maciousService.Delete(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
func (b *BaseApi) SearchSandfly(c *gin.Context) {
	var req dto.SearchSandflyWithPage
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	total, list, err := maciousService.SearchWithPage(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}
func (b *BaseApi) SearchSandflyRecord(c *gin.Context) {
	var req dto.SandflyRecordSearch
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	total, list, err := maciousService.LoadSandflyRecords(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// @Tags Macious
// @Summary Clean sandfly record
// @Description 清空扫描报告
// @Accept json
// @Param request body dto.OperateByID true "request"
// @Security ApiKeyAuth
// @Router /toolbox/macious/record/clean [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":true,"db":"sandflys","output_column":"name","output_value":"name"}],"formatZH":"清空恶意文件扫描报告 [name]","formatEN":"clean sandfly record [name]"}
func (b *BaseApi) CleanSandflyRecord(c *gin.Context) {
	var req dto.OperateByID
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := maciousService.CleanRecord(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}
