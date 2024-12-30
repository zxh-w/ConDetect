package v1

import (
	"ConDetect/backend/app/api/v1/helper"
	"ConDetect/backend/app/dto"
	"ConDetect/backend/constant"

	"github.com/gin-gonic/gin"
)

func (b *BaseApi) SearchMeasureRecord(c *gin.Context) {
	var req dto.SearchMeasureWithPage
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	total, list, err := trustedServie.SearchWithPage(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// @Tags Trusted
// @Summary Measure container credibility
// @Description 容器可信度量
// @Accept json
// @Param request body dto.OperationWithName true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /toolbox/measure/container [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"度量容器 [name]","formatEN":"measure container [name]"}
func (b *BaseApi) HandleContainerMeasure(c *gin.Context) {
	var req dto.OperationWithName
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := trustedServie.MeasureContainer(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
