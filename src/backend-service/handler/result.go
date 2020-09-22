package handler

import (
	m "app-deploy-platform/backend-service/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetResult(c *gin.Context) {
	result := m.NewResult()

	deployEnv := c.DefaultQuery("deployEnv", "prd")
	name := c.DefaultQuery("name", "test")
	db := m.Model
	db = db.Where("name = ?", name)
	db = db.Where("deploy_env = ?", deployEnv)
	db = db.Order("created_at desc")
	db.Limit(1).Find(result)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
		"data": result,
	})
}

func PostFormResult(c *gin.Context) {
	result := m.NewResult()

	result.Name = c.PostForm("name")
	result.DeployEnv = c.PostForm("deploy_env")
	result.Version = c.PostForm("version")

	m.Model.Create(result)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
	})
}

func PostResult(c *gin.Context) {
	result := m.NewResult()

	if err := c.ShouldBindJSON(&result); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "fail",
			"res":  "fail",
		})
		return
	}

	m.Model.Create(result)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
	})
}
