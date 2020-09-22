package config

import (
	"app-deploy-platform/common/tools"
)

type Env struct {
	Debug                bool
	AppName              string
	ServerPort           string
	DockerHarborUrl      string
	DockerHarborUser     string
	DockerHarborPassword string
}

func GetEnv() *Env {
	return &env
}

var env = Env{
	AppName:              tools.GetEnvDefault("APP_NAME", "service-operate-docker-harbor").(string),
	ServerPort:           tools.GetEnvDefault("SERVER_PORT", "80").(string),
	DockerHarborUrl:      tools.GetEnvDefault("DOCKER_HARBOR_URL", "https://docker.dm-ai.cn").(string),
	DockerHarborUser:     tools.GetEnvDefault("DOCKER_HARBOR_USER", "admin").(string),
	DockerHarborPassword: tools.GetEnvDefault("DOCKER_HARBOR_USER", "B26e2663X873").(string),
	Debug:                tools.GetEnvDefault("DEBUG_MODEL", true).(bool),
}
