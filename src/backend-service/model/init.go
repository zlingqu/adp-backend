package model

import (
	"app-deploy-platform/backend-service/config"
	"github.com/jinzhu/gorm"
	log "github.com/zuoshenglo/libs/logs/logrus"
)

var Model *gorm.DB

func init() {
	var err error
	log.Info(config.GetEnv().Database.FormatDSN())
	Model, err = gorm.Open("mysql", config.GetEnv().Database.FormatDSN())
	Model.LogMode(true)

	if err != nil {
		panic(err)
	}
}
