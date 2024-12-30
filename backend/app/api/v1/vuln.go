package v1

import (
	"ConDetect/backend/app/api/v1/helper"
	"ConDetect/backend/app/dto"
	"ConDetect/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags Vuln
// @Summary Handle trivy scan
// @Description 执行漏洞扫描
// @Accept json
// @Param request body dto.OperateByID true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/vuln/trivy/handle [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":true,"db":"trivies","output_column":"name","output_value":"name"}],"formatZH":"执行漏洞扫描 [name]","formatEN":"handle trivy scan [name]"}
func (b *BaseApi) HandleTrivyScan(c *gin.Context) {
	var req dto.OperateByID
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := vulnService.HandleOnce(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Vuln
// @Summary Create trivy
// @Description 创建扫描规则
// @Accept json
// @Param request body dto.TrivyCreate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/vuln/trivy/create [post]
// @x-panel-log {"bodyKeys":["name","target"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建漏洞扫描规则 [name][target]","formatEN":"create trivy [name][target]"}
func (b *BaseApi) CreateTrivy(c *gin.Context) {
	var req dto.TrivyCreate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := vulnService.Create(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Vuln
// @Summary Update trivy
// @Description 修改扫描规则
// @Accept json
// @Param request body dto.TrivyUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/vuln/trivy/update [post]
// @x-panel-log {"bodyKeys":["name","target"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"修改漏洞扫描规则 [name][target]","formatEN":"update trivy [name][target]"}
func (b *BaseApi) UpdateTrivy(c *gin.Context) {
	var req dto.TrivyUpdate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := vulnService.Update(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Vuln
// @Summary Delete trivy
// @Description 删除扫描规则
// @Accept json
// @Param request body dto.TrivyDelete true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/vuln/trivy/del [post]
// @x-panel-log {"bodyKeys":["ids"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"ids","isList":true,"db":"trivies","output_column":"name","output_value":"names"}],"formatZH":"删除漏洞扫描规则 [names]","formatEN":"delete trivy [names]"}
func (b *BaseApi) DeleteTrivy(c *gin.Context) {
	var req dto.TrivyDelete
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := vulnService.Delete(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
func (b *BaseApi) SearchTrivy(c *gin.Context) {
	var req dto.SearchTrivyWithPage
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	total, list, err := vulnService.SearchWithPage(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}
func (b *BaseApi) SearchTrivyRecord(c *gin.Context) {
	var req dto.TrivyRecordSearch
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	total, list, err := vulnService.LoadTrivyRecords(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// @Tags Vuln
// @Summary Clean trivy record
// @Description 清空扫描报告
// @Accept json
// @Param request body dto.OperateByID true "request"
// @Security ApiKeyAuth
// @Router /toolbox/vuln/record/clean [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":true,"db":"trivies","output_column":"name","output_value":"name"}],"formatZH":"清空漏洞扫描报告 [name]","formatEN":"clean trivy record [name]"}
func (b *BaseApi) CleanTrivyRecord(c *gin.Context) {
	var req dto.OperateByID
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := vulnService.CleanRecord(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}
