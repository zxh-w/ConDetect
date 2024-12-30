package v1

import (
	"ConDetect/backend/app/api/v1/helper"
	"ConDetect/backend/app/dto"
	"ConDetect/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags Host_User
// @Summary Update User
// @Description 修改用户权限
// @Accept json
// @Param request body dto.UserMmanage true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /hosts/user/update [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"修改主机用户权限 [name]","formatEN":"update hosts user [name]"}
func (b *BaseApi) UpdateUserManage(c *gin.Context) {
	var req dto.UserMmanage
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := userManageService.Update(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Host_User
// @Summary Delete User
// @Description 删除主机用户
// @Accept json
// @Param request body dto.SandflyDelete true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /hosts/user/delete [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"删除主机用户 [name]","formatEN":"delete host user [name]"}
func (b *BaseApi) DeleteUserManage(c *gin.Context) {
	var req dto.UserMmanage
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := userManageService.Delete(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

func (b *BaseApi) SearchUserManage(c *gin.Context) {
	var req dto.UserSearch
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	_, list, err := userManageService.SearchWithName(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, list)
}
