package middleware

import (
	"errors"
	"strings"

	"ConDetect/backend/app/api/v1/helper"
	"ConDetect/backend/app/repo"
	"ConDetect/backend/constant"

	"github.com/gin-gonic/gin"
)

func BindDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		settingRepo := repo.NewISettingRepo()
		status, err := settingRepo.Get(settingRepo.WithByKey("BindDomain"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}
		if len(status.Value) == 0 {
			c.Next()
			return
		}
		domains := c.Request.Host
		parts := strings.Split(c.Request.Host, ":")
		if len(parts) > 0 {
			domains = parts[0]
		}

		if domains != status.Value {
			if LoadErrCode("err-domain") != 200 {
				helper.ErrResponse(c, LoadErrCode("err-domain"))
				return
			}
			helper.ErrorWithDetail(c, constant.CodeErrDomain, constant.ErrTypeInternalServer, errors.New("domain not allowed"))
			return
		}
		c.Next()
	}
}
