package router

import (
	"app-deploy-platform/backend-service/config"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	log "github.com/zuoshenglo/libs/logs/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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
	// request log
	engine.Use(handRequestLog())
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

func handleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				log.Error(err)

				var (
					errMsg     string
					mysqlError *mysql.MySQLError
					ok         bool
				)
				if errMsg, ok = err.(string); ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "system error, " + errMsg,
					})
					return
				} else if mysqlError, ok = err.(*mysql.MySQLError); ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "system error, " + mysqlError.Error(),
					})
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "system error",
					})
					return
				}
			}
		}()
		c.Next()
	}
}

// request log
func handRequestLog() gin.HandlerFunc {
	var reqLog = logrus.New()
	reqLog.Out = os.Stdout
	reqLog.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "date",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "@caller",
		},
	}

	var requestLog = log.NewRequestLog()

	return func(c *gin.Context) {
		startTime := time.Now()
		// set requestLog
		requestLog.Request = c.Request.RequestURI
		requestLog.RemoteAddr = c.Request.RemoteAddr
		body, _ := ioutil.ReadAll(c.Request.Body)
		//requestLog.RequestBody = c.Request.Body.Read()
		//requestLog.Status = c.Request.Response.Status
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		c.Next()
		endTime := time.Now()
		// end set requestLog
		requestLog.StartTime = startTime
		requestLog.EndTime = endTime
		requestLog.Status = c.Writer.Status()
		requestLog.RequestTime = endTime.Sub(startTime)
		requestLog.RequestUrl = c.Request.RequestURI
		requestLog.ReqMethod = c.Request.Method
		requestLog.ClientIP = c.ClientIP()
		requestLog.RequestBody = string(body)
		reqLog.WithFields(logrus.Fields{
			"type":    "pa_access",
			"channel": config.GetEnv().AppName,
			"msg":     requestLog,
		}).Info("request info")
	}
}

func HealCheck(c *gin.Context) {
	c.String(200, "# TYPE health_info gauge\n"+"health_info{status=\"ok\", name=\"service-adp-deploy\", namespace=\"devops\"} 0\n")
}
