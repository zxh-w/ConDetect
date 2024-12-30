package middleware

import (
	"strconv"

	"ConDetect/backend/app/api/v1/helper"
	"ConDetect/backend/app/repo"
	"ConDetect/backend/constant"
	"ConDetect/backend/global"

	"github.com/gin-gonic/gin"
)

func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if method, exist := c.Get("authMethod"); exist && method == constant.AuthMethodJWT {
			c.Next()
			return
		}
		sId, err := c.Cookie(constant.SessionName)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrUnauthorized, constant.ErrTypeNotLogin, nil)
			return
		}
		psession, err := global.SESSION.Get(sId)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrUnauthorized, constant.ErrTypeNotLogin, nil)
			return
		}
		settingRepo := repo.NewISettingRepo()
		setting, err := settingRepo.Get(settingRepo.WithByKey("SessionTimeout"))
		if err != nil {
			global.LOG.Errorf("create operation record failed, err: %v", err)
		}
		lifeTime, _ := strconv.Atoi(setting.Value)
		_ = global.SESSION.Set(sId, psession, lifeTime)
		c.Next()
	}
}
