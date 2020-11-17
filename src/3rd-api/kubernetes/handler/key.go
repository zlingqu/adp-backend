package handler

import (
	"app-deploy-platform/3rd-api/kubernetes/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetK8sKeyFile(c *gin.Context) {
	env := c.Query("env")

	if env == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  "请确认接口参数的准确。",
		})
		return
	}

	downFile := config.GetKeyFile(env)
	if downFile == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  "参数错误，找不到指定的key文件",
		})
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "config"))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(downFile)
}
