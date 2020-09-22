package config

import (
	"app-deploy-platform/common/tools"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func GetConfFile(fileName string) string {
	log.Info(fmt.Sprintf("查找配置文件%s", fileName))
	dir, _ := os.Getwd()
	confFile := dir + "/key/" + fileName
	log.Info(fmt.Sprintf("配置文件为：%s", confFile))
	return confFile
}

type Env struct {
	Debug      bool
	AppName    string
	ServerPort string
}

func GetEnv() *Env {
	return &env
}

var env = Env{
	AppName:    tools.GetEnvDefault("APP_NAME", "service-k8s-app-status-check-v1").(string),
	ServerPort: tools.GetEnvDefault("SERVER_PORT", "80").(string),
	Debug:      tools.GetEnvDefault("DEBUG_MODEL", true).(bool),
}
