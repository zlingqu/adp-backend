package handler

import (
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/common/tools"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetEnv(c *gin.Context) {

	var env []m.Env
	var getEnv m.GetEnv
	var count int64
	if err := c.ShouldBind(&getEnv); err != nil {
		log.Error(err)
	}

	m.Model.Where("name LIKE ?", "%"+getEnv.Name+"%").Find(&env).Count(&count)
	log.Println(env)
	log.Println(count)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"msg":   "ok",
		"data":  env,
	})
}

func GetEnvByID(c *gin.Context) {

	var env m.Env
	var count int64
	var getEnvByID m.GetEnvByID
	if err := c.ShouldBindUri(&getEnvByID); err != nil {
		log.Error(err)
		// return
	}

	log.Println(getEnvByID.ID)
	RowsAffected := m.Model.First(&env, getEnvByID.ID).Count(&count).RowsAffected
	if RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  404,
			"count": 0,
			"msg":   fmt.Sprintf("找不到id=%s的env", getEnvByID.ID),
			"data":  nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": 1,
		"msg":   "ok",
		"data":  env,
	})
}

func PostEnv(c *gin.Context) {
	env := m.NewEnv()
	if err := c.ShouldBindJSON(&env); err != nil {
		log.Error(err)
		// return
	}
	m.Model.Create(env)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func PostEnvs(c *gin.Context) {
	var env []m.Env
	var postEnvs m.PostEnvs
	var count int64
	if err := c.ShouldBind(&postEnvs); err != nil {
		log.Error(err)
		// return
	}

	m.Model.Where("id in (?)", postEnvs.Ids).Find(&env).Count(&count)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"msg":   "ok",
		"data":  env,
	})
}

func PutEnv(c *gin.Context) {
	env := m.NewEnv()
	if err := c.ShouldBind(env); err != nil {
		log.Error(err)
	}

	log.Println(*env)

	//Model.Save(env)
	m.Model.Model(env).Updates(map[string]interface{}{"name": env.Name, "status": env.Status})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func DeleteEnv(c *gin.Context) {
	env := m.NewEnv()
	env.ID = tools.StringToUint(c.Param("id"))
	m.Model.Delete(env)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
	})

}
