package v1

import (
	"ConDetect/backend/app/api/v1/helper"
	"ConDetect/backend/app/dto"
	"ConDetect/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags Baseline
// @Summary Handle baselin check
// @Description 执行基线核查
// @Accept json
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/baseline/check [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFunctions":[],"formatZH":"执行基线核查","formatEN":"perform baseline check"}
func (b *BaseApi) CheckBaseline(c *gin.Context) {
	if err := baselineService.CheckBaseline(); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

func (b *BaseApi) SearchBaseline(c *gin.Context) {
	var req dto.SearchBaselineWithPage
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	total, datas, err := baselineService.SearchBaseline(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: datas,
		Total: total,
	})
}
