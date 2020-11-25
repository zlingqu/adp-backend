package model

import (
	"app-deploy-platform/backend-service/config"

	"gorm.io/driver/mysql"

	// "github.com/jinzhu/gorm"
	log "github.com/zuoshenglo/libs/logs/logrus"
	"gorm.io/gorm"
)

// Model 定义db示例
var Model gorm.Migrator
var DB *gorm.DB

func init() {
	var err error
	log.Info(config.GetEnv().Database.FormatDSN())
	// Model, err = gorm.Open("mysql", config.GetEnv().Database.FormatDSN())
	DB, err = gorm.Open(mysql.Open(config.GetEnv().Database.FormatDSN()), &gorm.Config{})

	Model = DB.Migrator()
	// Model.LogMode(true)

	if err != nil {
		panic(err)
	}
}
