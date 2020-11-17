package handler

import (
	m "app-deploy-platform/backend-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetResult(c *gin.Context) {
	// var result []m.Result
	result := m.NewResult()

	deployEnv := c.DefaultQuery("deployEnv", "prd")
	name := c.DefaultQuery("name", "test")
	db := m.Model
	db.Where("name = ?", name).Where("deploy_env = ?", deployEnv).Order("created_at desc").Limit(1).Find(result)
	// if len(result) == 0 {
	if result.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 1,
			"msg":  "error",
			"res":  "error",
			"data": "查询不到数据",
		})
		return
	}
	// db = db.Where("deploy_env = ?", deployEnv)
	// db = db.Order("created_at desc")
	// db.Limit(3).Find(result)
	// fmt.Printf("%#v", result)
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
