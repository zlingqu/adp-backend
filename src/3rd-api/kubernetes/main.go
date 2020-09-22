package main

import (
	"app-deploy-platform/3rd-api/kubernetes/config"
	"app-deploy-platform/3rd-api/kubernetes/router"
	"app-deploy-platform/common/server"
)

func main() {
	server.Run(router.InitEngine(), config.GetEnv().ServerPort, config.GetEnv().Debug)
}
