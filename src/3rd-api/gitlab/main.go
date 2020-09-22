package main

import (
	"app-deploy-platform/3rd-api/gitlab/config"
	"app-deploy-platform/3rd-api/gitlab/router"
	"app-deploy-platform/common/server"
)

func main() {
	server.Run(router.InitEngine(), config.GetEnv().ServerPort, config.GetEnv().Debug)
}
