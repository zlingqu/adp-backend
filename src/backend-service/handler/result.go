package handler

import (
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/backend-service/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func PostFormResult(c *gin.Context) (int, interface{}, string) {
	result := m.NewResult()

	result.Name = c.PostForm("name")
	result.DeployEnv = c.PostForm("deploy_env")
	result.Version = c.PostForm("version")

	m.DB.Create(result)

	return 0, "ok", "ok"
}

func PostResult(c *gin.Context) (int, interface{}, string) {
	result := m.NewResult()

	if err := c.ShouldBindJSON(&result); err != nil {
		log.Error(err)
		return 0, "fail", "fail"
	}

	m.DB.Create(result)
	server.PushResult(*result)

	return 0, "ok", "ok"
}
