package router

import (
	"app-deploy-platform/3rd-api/kubernetes/handler"
	r "app-deploy-platform/common/router"
	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) *gin.Engine {
	return r.Init(engine, func(rgs map[r.Api]*gin.RouterGroup) {
		v1 := rgs[r.ApiV1]
		v1.GET("/pods-status", handler.GetPodsStatus)
		v1.GET("/get-k8s-key-file", handler.GetK8sKeyFile)
	})
}

func InitEngine() *gin.Engine {
	return Init(gin.New())
}
