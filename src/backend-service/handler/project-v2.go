package handler

import (
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/common/tools"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetProjectV2IdNameGitLang(c *gin.Context) {
	projectName := c.DefaultQuery("name", "")
	db := m.Model
	var repData []m.ProjectIdNameGitLang

	// select db
	db = db.Table(m.Project{}.TableName())
	db = db.Select("id, name, git_repository, language_type, if_use_model, if_use_git_manager_model, model_git_repository")
	db = db.Where("name like ?", "%"+projectName+"%")
	db.Scan(&repData)

	// reponse
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
		"data": repData,
	})
}

func GetProjectV2IdNameGitLangProduct(c *gin.Context) {
	projectName := c.DefaultQuery("name", "")
	db := m.Model
	var repData []m.ProjectIdNameGitLangProduct

	// select db
	db = db.Table(m.Project{}.TableName())
	db = db.Select("id, name, git_repository, language_type, owned_product, if_use_model, if_use_git_manager_model, model_git_repository")
	db = db.Where("name like ?", "%"+projectName+"%")
	db.Scan(&repData)

	// reponse
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
		"data": repData,
	})
}

func GetProjectV2IdNameGit(c *gin.Context) {
	projectName := c.DefaultQuery("name", "")
	db := m.Model
	var repData []m.ProjectIdNameGit

	// select db
	db = db.Table(m.Project{}.TableName())
	db = db.Select("id, name, git_repository")
	db = db.Where("name like ?", "%"+projectName+"%")
	db.Scan(&repData)

	// reponse
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
		"data": repData,
	})
}

func GetProjectV2IdName(c *gin.Context) {
	projectName := c.DefaultQuery("name", "")
	languageType := c.DefaultQuery("language_type", "")

	db := m.Model
	var repData []m.ProjectIdName

	// select db
	db = db.Table(m.Project{}.TableName())
	db = db.Select("id, name")
	db = db.Where("name like ?", "%"+projectName+"%")
	if languageType != "" {
		db = db.Where("language_type = ?", languageType)
	}
	db.Scan(&repData)

	// reponse
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
		"data": repData,
	})
}

func GetProjectV2(c *gin.Context) {

	var project []m.Project
	//var getProject m.GetProject
	var count int64

	// req param
	name := strings.TrimSpace(c.DefaultQuery("name", ""))
	ownedProduct := strings.TrimSpace(c.DefaultQuery("owned_product", ""))
	deployEnvType := strings.TrimSpace(c.DefaultQuery("deploy_env_type", ""))
	languageType := strings.TrimSpace(c.DefaultQuery("language_type", ""))

	limit, offset := tools.GetMysqlLimitOffset(c.DefaultQuery("page", "1"), c.DefaultQuery("size", "10"))
	log.Println(fmt.Sprintf("req parms: name: %s, limit: %d, offset: %d", name, limit, offset))
	db := m.Model.Where("name LIKE ?", "%"+name+"%")

	if ownedProduct != "" {
		db = db.Where("owned_product in (?)", strings.Split(ownedProduct, ","))
	}

	if deployEnvType != "" {
		db = db.Where("deploy_env_type in (?)", strings.Split(deployEnvType, ","))
	}

	if languageType != "" {
		db = db.Where("language_type in (?)", strings.Split(languageType, ","))
	}
	log.Println(ownedProduct, deployEnvType, languageType)
	log.Println(len(ownedProduct), len(deployEnvType), len(languageType))

	db.Limit(limit).Offset(offset).Find(&project)
	db.Model(&m.Project{}).Count(&count)
	//m.Model.Limit(limit).Offset(offset).Where("name LIKE ?", "%" + name + "%").Find(&project)
	//m.Model.Model(&m.Project{}).Where("name LIKE ?", "%" + name + "%").Count(&count)

	log.Println(project)
	log.Println(count)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"msg":   "ok",
		"res":   "ok",
		"data":  project,
	})
}
