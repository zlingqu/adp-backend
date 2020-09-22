package main

import (
	"app-deploy-platform/3rd-api/harbor/config"
	"app-deploy-platform/3rd-api/harbor/router"
	"app-deploy-platform/common/server"
)

func main() {
	server.Run(router.InitEngine(), config.GetEnv().ServerPort, config.GetEnv().Debug)
}
