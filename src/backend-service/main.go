package main

import (
	"app-deploy-platform/backend-service/config"
	"app-deploy-platform/backend-service/database"
	"app-deploy-platform/backend-service/router"
	"app-deploy-platform/backend-service/server"
	"encoding/json"
	"fmt"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	env := config.GetEnv()
	marshal, _ := json.Marshal(env)
	fmt.Println(string(marshal))

	server.RunEventSource(env.EventSourcePort)

	database.Init(env.Database)

	apiServer := server.RunApi(router.InitEngine(), env.ServerPort, env.Debug)

	server.WaitInterrupt(apiServer)

}
