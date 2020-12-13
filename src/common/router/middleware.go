package router

import (
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

type HandlerFunc func(*gin.Context) (code int, data interface{}, msg string)

func WrapHandlerFunc(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, data, msg := handlerFunc(c)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"res":  data,
			"msg":  msg,
		})
	}
}

// 错误统一处理
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

// 打印请求日志
func handRequestLog(appName string) gin.HandlerFunc {
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
			"channel": appName,
			"msg":     requestLog,
		}).Info("request info")
	}
}
