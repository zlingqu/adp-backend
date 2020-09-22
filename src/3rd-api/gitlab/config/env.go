package config

import (
	"app-deploy-platform/common/tools"
)

type Env struct {
	Debug                          bool
	ServerPort                     string
	PrivateToken                   string
	GitLabOperateProjectApiAddress string
}

func GetEnv() *Env {
	return &env
}

// 本文件建议在代码协同工具(git/svn等)中忽略

var env = Env{
	Debug:                          tools.GetEnvDefault("DEBUG_MODEL", true).(bool),
	ServerPort:                     tools.GetEnvDefault("SERVER_PORT", "80").(string),
	PrivateToken:                   "fr_k4PoP95fCxs8AoQzx",
	GitLabOperateProjectApiAddress: "https://gitlab.dm-ai.cn/api/v4/projects",
}
