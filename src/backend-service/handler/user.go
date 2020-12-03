package handler

import (
	m "app-deploy-platform/backend-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserForName(c *gin.Context) {
	// init all param
	userName := c.DefaultQuery("name", "")
	db := m.DB
	var repData []m.UserInfo

	// select db
	db = db.Table("user")
	db = db.Select("owner_english_name, owner_china_name")
	db = db.Where("owner_english_name like ?", "%"+userName+"%").Or("owner_china_name like ?", "%"+userName+"%")
	db.Scan(&repData)

	// reponse
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
		"data": repData,
	})
}

func GetUserChinaName(c *gin.Context) { //修改工单接口，会对用户做判定
	ownerEnglishName := c.DefaultQuery("ownerEnglishName", "")
	db := m.DB
	var repData struct {
		OwnerChinaName string `json:"owner_china_name"`
	}

	// select db
	db = db.Table("user")
	db = db.Select("owner_china_name")
	db = db.Where("owner_english_name = ?", ownerEnglishName)
	db.Find(&repData)

	// reponse
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
		"data": repData,
	})
}

func SyncLdapUser(c *gin.Context) {

	db := m.DB

	misLdapService := m.NewMisLdapService()
	misLdapService.Request()

	// 处理请求的结果
	if misLdapService.Res != "ok" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"res":  misLdapService.Res,
			"msg":  misLdapService.Msg,
		})
		return
	}
	// 然后同步 请求的数据到 mysql指定的数据表中。
	for _, v := range misLdapService.Rep.Data {
		user := m.User{
			OwnerChinaName:   v.DisplayName,
			OwnerEnglishName: v.Cn,
		}
		db.Where(m.User{OwnerEnglishName: v.Cn}).FirstOrCreate(&user)

	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "ok",
		"msg":  "同步更新用户数据成功。",
	})
}
