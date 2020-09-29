package handler

import (
	conf "app-deploy-platform/3rd-api/kubernetes/config"
	m "app-deploy-platform/3rd-api/kubernetes/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPodsStatus(c *gin.Context) {

	// init
	env := c.DefaultQuery("env", "")
	namespace := c.DefaultQuery("namespace", "")
	appName := c.DefaultQuery("appName", "")
	imageSha := c.DefaultQuery("imageSha", "")

	// check request param
	if env == "" || namespace == "" || appName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  "请确认接口参数的准确。",
		})
	}

	//pods status
	pt := m.NewPodsStatus()
	// pt.SetKubernetesConfigFilePath(env)
	pt.KubernetesConfigFile = conf.GetKeyFile(env)
	if pt.KubernetesConfigFile == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  "env参数值错误，找不到对应的key文件",
		})
		return
	}

	pt.SetKubernetesClient()
	pt.GetPodsInfo(namespace, appName, imageSha)

	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"msg":    pt.Res.Msg,
		"res":    pt.Res.Res,
		"status": pt.Res.Status,
	})
}
