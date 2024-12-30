package service

import (
	"ConDetect/backend/constant"
	"ConDetect/backend/global"
	"context"

	"gorm.io/gorm"
)

func getTxAndContext() (tx *gorm.DB, ctx context.Context) {
	tx = global.DB.Begin()
	ctx = context.WithValue(context.Background(), constant.DB, tx)
	return
}
