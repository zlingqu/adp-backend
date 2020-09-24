package handler

import (
	"app-deploy-platform/3rd-api/harbor/config"
	"app-deploy-platform/3rd-api/harbor/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DockerImageSha256(c *gin.Context) {

	space := c.DefaultQuery("space", "")
	project := c.DefaultQuery("project", "")
	tag := c.DefaultQuery("tag", "")
	dh := model.NewDockerHarbor(space, project, config.GetEnv().DockerHarborUser, config.GetEnv().DockerHarborPassword)

	c.JSON(http.StatusOK, gin.H{
		"data": dh.GetDigest(tag),
		"code": 0,
		"res":  dh.Res,
		"msg":  dh.Msg,
	})

}
