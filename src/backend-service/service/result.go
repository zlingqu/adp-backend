package service

import (
	"app-deploy-platform/backend-service/model"
	"time"
)

func GetResultByCreateTime(t time.Time) (results []model.Result) {
	model.DB.Where("created_at >= ?", t).Find(&results)
	return
}
