package config

import (
	"app-deploy-platform/common/tools"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func GetKeyFile(envName string) string {
	log.Info(fmt.Sprintf("查找配置文件%s", envName))
	dir, _ := os.Getwd()
	path := dir + "/key/" + envName
	if _, err := os.Stat(path); err != nil {

		log.Info(fmt.Sprintf("%s不存在\n", path))
		return ""

	}
	log.Info(fmt.Sprintf("path %s 存在\n", path))
	return path

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
