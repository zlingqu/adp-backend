package main

import (
	// "fmt"
	
	"app-deploy-platform/3rd-api/gitlab/config"
	"app-deploy-platform/3rd-api/gitlab/router"
	"app-deploy-platform/common/server"
	// svc "app-deploy-platform/3rd-api/gitlab/service"
)

func main() {
	server.Run(router.InitEngine(), config.GetEnv().ServerPort, config.GetEnv().Debug)
	// id := svc.GetCommitIDByRepourlAndBranch("https://gitlab.dm-ai.cn/liaobin/dm-svc-component-boilerplate.git", "dev")
	// fmt.Println(id)
}
