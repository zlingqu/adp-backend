package router

import (
	. "app-deploy-platform/backend-service/handler"
	r "app-deploy-platform/common/router"

	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) *gin.Engine {
	return r.Init(engine, func(rgs map[r.Api]*gin.RouterGroup) {
		v1 := rgs[r.ApiV1]
		v2 := rgs[r.ApiV2]

		deployRoute(v1, v2)
		envRoute(v1, v2)
		projectRoute(v1, v2)
		spaceRoute(v1, v2)
		userRoute(v1, v2)
		tools(v1, v2)

		v1.POST("/result", r.WrapHandlerFunc(PostResult))
		v1.POST("/result-form", r.WrapHandlerFunc(PostFormResult))
	})

}

func InitEngine() *gin.Engine {
	return Init(gin.New())
}

func deployRoute(v1 *gin.RouterGroup, v2 *gin.RouterGroup) {
	v1.GET("/healthCheck", HealCheck)
	v1.GET("/metrics", HealCheck)
	v1.GET("/deploy/online/:id", DeployOnline)
	v1.POST("/deploy/onlines", Mdeploy)
	v1.POST("/deployments/list", PostDeployList)
	v1.POST("/deployments/create", PostDeploy)
	v1.POST("/deployments/update", PostUpdate)
	v1.POST("/deployments/change", PostChange)
	v1.DELETE("/deploy/:id", DeleteDeploy)
	v2.GET("/deploy", GetDeploy)
}

func envRoute(v1 *gin.RouterGroup, v2 *gin.RouterGroup) {
	v1.GET("/env", GetEnv)
	v1.GET("/env/:id", GetEnvByID)
	v1.POST("/env", PostEnv)
	v1.POST("/envs", PostEnvs)
	v1.PUT("/env", PutEnv)
	v1.DELETE("/env/:id", DeleteEnv)
}

func projectRoute(v1 *gin.RouterGroup, v2 *gin.RouterGroup) {
	v1.GET("/project", GetProject)         //http://{{host}}/api/v1/project/?name=backend，如果没有name选项表示查看所有
	v1.GET("/project/:id", GetProjectById) //http://{{host}}/api/v1/project/80
	v1.POST("/project", PostProject)
	v1.POST("/projects", PostProjects)
	v1.PUT("/project", PutProject)
	v1.DELETE("/project/:id", DeleteProject)

	v2.GET("/project", GetProjectV2)
	v2.GET("/project-id-name", GetProjectV2IdName)
	v2.GET("/project-id-name-git", GetProjectV2IdNameGit)
	v2.GET("/project-id-name-git-lang", GetProjectV2IdNameGitLang)
	v2.GET("/project-id-name-git-lang-product", GetProjectV2IdNameGitLangProduct)
}

func spaceRoute(v1 *gin.RouterGroup, v2 *gin.RouterGroup) {
	v1.GET("/space", GetSpace)
	v1.POST("/space", PostSpace)
	v1.PUT("/space/:id", PutSpace)
	v1.DELETE("/space/:id", DeleteSpace)
}

func userRoute(v1 *gin.RouterGroup, v2 *gin.RouterGroup) {
	v1.GET("/sync-ldap-user", SyncLdapUser)
	v1.GET("/get-user-for-name", GetUserForName)
	v1.GET("/user/get-owner-china-name", GetUserChinaName)
}

func tools(v1 *gin.RouterGroup, v2 *gin.RouterGroup) {
	v1.GET("/tools/qrcode",GetQrcode)

}

