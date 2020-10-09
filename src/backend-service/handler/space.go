package handler

import (
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/common/tools"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//GetSpace 获取space列表，/api/v1/space?name=xmc&page=2&size=10，其中name可以为空。page和size两个参数没有使用到
func GetSpace(c *gin.Context) {
	var space []m.Space
	var param m.GetSpace
	var count int64
	if err := c.ShouldBind(&param); err != nil {
		log.Error(err)
	}

	m.Model.Where("name LIKE ?", "%"+param.Name+"%").Find(&space).Count(&count)
	log.Info("GetSpace查出条数", count)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"msg":   "ok",
		"data":  space,
	})
}

// PostSpace 创建space，POST /api/v1/space body是json格式
func PostSpace(c *gin.Context) {
	space := m.NewSpace()
	if err := c.ShouldBindJSON(&space); err != nil {
		log.Error(err)
	}
	log.Info(space)
	rows := m.Model.Create(space).RowsAffected
	if rows == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "插入数据库失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

// PutSpace 用于更新space,/api/v1/space/48.请求体是json,{name: "devops", owner: "quzhongling"}
func PutSpace(c *gin.Context) {
	space := m.NewSpace()

	if err := c.ShouldBind(space); err != nil {
		log.Error(err)
	}

	id := c.Param("id")
	space.ID = tools.StringToUint(id)

	log.Println(*space)

	raws := m.Model.Model(space).Updates(*space).RowsAffected
	if raws == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  fmt.Sprintf("更新失败，表%s中没有这样的记录id=%s", space.TableName(), id),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

// DeleteSpace 删除space，根据id 。/api/v1/space/116
func DeleteSpace(c *gin.Context) {
	space := m.NewSpace()
	id := c.Param("id")

	space.ID = tools.StringToUint(id)
	log.Println("del id : ", id)

	RowsAffected := m.Model.Delete(space).RowsAffected
	if RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  fmt.Sprintf("请求错误，没有id=%s这样的space", id),
			"res":  "error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
	})

}
