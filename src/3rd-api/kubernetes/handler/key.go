package handler

import (
	"app-deploy-platform/3rd-api/kubernetes/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetK8sKeyFile(c *gin.Context) {
	env := c.DefaultQuery("env", "")

	if env == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  "请确认接口参数的准确。",
		})
		return
	}

	// 获得指定的配置文件
	downFile := config.GetConfFile(env)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "config"))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(downFile)
}
