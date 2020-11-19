package main

// import (
// 	svc "app-deploy-platform/3rd-api/gitlab/service"
// 	"fmt"
// )

// func main() {
// 	id := svc.GetTagsByRepourl("https://gitlab.dm-ai.cn/X4C/AI_GZ/Android.git")
// 	fmt.Println(id)
// }

import (
	"app-deploy-platform/3rd-api/gitlab/config"
	"app-deploy-platform/3rd-api/gitlab/router"
	"app-deploy-platform/common/server"
)

func main() {
	server.Run(router.InitEngine(), config.GetEnv().ServerPort, config.GetEnv().Debug)

}
