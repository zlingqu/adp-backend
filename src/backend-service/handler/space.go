package handler

import (
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/common/tools"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetSpace(c *gin.Context) {

	var space []m.Space
	var getSpace m.GetSpace
	var count int64
	if err := c.ShouldBind(&getSpace); err != nil {
		log.Error(err)
		// return
	}

	m.Model.Where("name LIKE ?", "%"+getSpace.Name+"%").Find(&space).Count(&count)
	log.Println(space)
	log.Println(count)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"msg":   "ok",
		"data":  space,
	})
}

func PostSpace(c *gin.Context) {
	space := m.NewSpace()

	if err := c.ShouldBindJSON(&space); err != nil {
		log.Error(err)
		// return
	}
	log.Println(space)
	m.Model.Create(space)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func PutSpace(c *gin.Context) {
	space := m.NewSpace()

	if err := c.ShouldBind(space); err != nil {
		log.Error(err)
	}

	id := c.Param("id")
	space.ID = tools.StringToUint(id)

	log.Println(*space)

	//m.Model.Save(env)
	m.Model.Model(space).Updates(*space)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func DeleteSpace(c *gin.Context) {
	space := m.NewSpace()
	id := c.Param("id")

	space.ID = tools.StringToUint(id)
	//check id , default 0
	log.Println("del id : ", id)
	//if err := c.ShouldBind(space); err != nil {
	//	log.Error(err)
	//}

	m.Model.Delete(space)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
	})

}
