package router

import (
	"app-deploy-platform/3rd-api/gitlab/handler"
	r "app-deploy-platform/common/router"

	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) *gin.Engine {
	return r.Init(engine, func(rgs map[r.Api]*gin.RouterGroup) {
		v1 := rgs[r.ApiV1]
		v1.GET("/gitlab/branchs", handler.GetBranchs)
		v1.GET("/gitlab/commit_id", handler.GetCommitID)
		v1.GET("/gitlab/tags", handler.GetTags)
		v1.GET("/gitlab/project", handler.ProjectInfo)
	})
}

func InitEngine() *gin.Engine {
	return Init(gin.New())
}
