package router

import (
	"app-deploy-platform/3rd-api/harbor/handler"
	r "app-deploy-platform/common/router"
	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) *gin.Engine {
	return r.Init(engine, func(rgs map[r.Api]*gin.RouterGroup) {
		v1 := rgs[r.ApiV1]
		v1.GET("/docker-image-sha256", handler.DockerImageSha256)
	})
}

func InitEngine() *gin.Engine {
	return Init(gin.New())
}
