package main

import (
	// "fmt"

	"app-deploy-platform/3rd-api/gitlab/config"
	"app-deploy-platform/3rd-api/gitlab/router"
	"app-deploy-platform/common/server"
	// svc "app-deploy-platform/3rd-api/gitlab/service"
	// "fmt"
)

func main() {
	server.Run(router.InitEngine(), config.GetEnv().ServerPort, config.GetEnv().Debug)
	// id := svc.GetTagsByRepourl("https://gitlab.dm-ai.cn/XMC/xmc-tk/xmc-offline-task.git")
	// fmt.Println(id)
}
