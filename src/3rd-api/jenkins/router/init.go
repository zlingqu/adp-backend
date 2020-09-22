package router

import (
	"app-deploy-platform/3rd-api/jenkins/handler"
	r "app-deploy-platform/common/router"
	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) *gin.Engine {
	return r.Init(engine, func(rgs map[r.Api]*gin.RouterGroup) {
		v1 := rgs[r.ApiV1]
		v1.POST("/multibranch-webhook-trigger", handler.PostMultibranchWebhookTrigger)
		v1.PUT("/job", handler.PutJob)
		v1.POST("/job", handler.PostJob)
		v1.DELETE("/job", handler.DeleteJob)
		v1.POST("/jenkins_job", handler.NewOperateJenkins)
		v1.POST("/jenkins/build", handler.NewJenkinsBuild)
	})
}

func InitEngine() *gin.Engine {
	return Init(gin.New())
}
