package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Api string

const (
	ApiV1 = "/api/v1"
	ApiV2 = "/api/v2"
)

var (
	isHealthCHeck = false
	routerGroups  map[Api]*gin.RouterGroup
)

func InitEngine(apiDefine func(routerGroups map[Api]*gin.RouterGroup)) *gin.Engine {
	return Init(gin.New(), apiDefine)
}

func Init(engine *gin.Engine, apiDefine func(routerGroups map[Api]*gin.RouterGroup)) *gin.Engine {

	if !isHealthCHeck {
		engine.GET("/healthCheck", HealCheck)
		isHealthCHeck = true
	}

	if routerGroups == nil {
		routerGroups = map[Api]*gin.RouterGroup{
			ApiV1: engine.Group(ApiV1),
			ApiV2: engine.Group(ApiV2),
		}
	}

	apiDefine(routerGroups)

	engine.Use(gin.Logger())

	engine.Use(handleErrors()) // 错误处理
	// TODO appName改成按模块取 现在改动的文件非常多
	engine.Use(handRequestLog("adp-backend"))
	//engine.Use(filters.RegisterSession()) // 全局session
	//engine.Use(filters.RegisterCache())   // 全局cache

	//engine.Use(auth.RegisterGlobalAuthDriver("cookie", "web_auth")) // 全局auth cookie
	//engine.Use(auth.RegisterGlobalAuthDriver("jwt", "jwt_auth"))    // 全局auth jwt

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
		return
	})

	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
		return
	})

	return engine

}

func HealCheck(c *gin.Context) {
	c.String(200, "# TYPE health_info gauge\n"+"health_info{status=\"ok\", name=\"service-adp-deploy\", namespace=\"devops\"} 0\n")
}
