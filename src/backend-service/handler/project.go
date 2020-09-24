package handler

import (
	"app-deploy-platform/backend-service/config"
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/common/tools"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"net/http"
	"strings"
)

func GetProject(c *gin.Context) {

	var project []m.Project
	var getProject m.GetProject
	var count int64
	if err := c.ShouldBind(&getProject); err != nil {
		log.Error(err)
		// return
	}

	m.Model.Where("name LIKE ?", "%"+getProject.Name+"%").Find(&project).Count(&count)
	log.Println(project)
	log.Println(count)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"msg":   "ok",
		"data":  project,
	})
}

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

func PostProject(c *gin.Context) {
	project := m.NewProject()
	var msg string
	var res string

	if err := c.ShouldBindJSON(&project); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "fail",
		})
		return
	}

	log.Println("The received front-end request data is : ", project,
		" , Start to request service-call-jenkins to create the project : ", project.Name)

	//msg = JenkinsProject(project, "create")
	//if msg != "ok" {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 0,
	//		"msg":  "fail",
	//	})
	//	return
	//}
	pj := NewProject()
	pj.AppName = project.Name
	pj.GitAddress = project.GitRepository
	res, msg = pj.CreateJob()
	if res != "ok" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  msg,
			"res":  res,
		})
		return
	}

	log.Println("Successfully created project in Jenkins，Start insert data in mysql-table project。")
	m.Model.Create(project)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
	})
}

func JenkinsProject(p *m.Project, action string) string {
	client := resty.New()

	j := m.JenkinsJob{
		GitAddress:  p.GitRepository,
		AppName:     p.Name,
		ProductName: p.OwnedProduct,
		Action:      action,
	}

	b, _ := json.Marshal(j)
	log.Println("The post data of the request service-operate-jenkins is:", string(b))

	r, e := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(b)).
		Post(config.GetEnv().JenkinsJobAddress)

	if e != nil {
		log.Error(e)
		return "fail"
	}

	log.Println("The response of service-operate-jenkins is:", r)
	var jenkinsResponse m.JenkinsResponse
	e = json.Unmarshal(r.Body(), &jenkinsResponse)

	if e != nil || jenkinsResponse.Status != "ok" {
		return "fail"
	}

	return "ok"
}

func PostProjects(c *gin.Context) {
	var project []m.Project
	var postIds m.PostIds
	var count int64
	if err := c.ShouldBind(&postIds); err != nil {
		log.Error(err)
		// return
	}

	m.Model.Where("id in (?)", postIds.Ids).Find(&project).Count(&count)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"msg":   "ok",
		"data":  project,
	})
}

func PutProject(c *gin.Context) {
	project := m.NewProject()
	if err := c.ShouldBind(project); err != nil {
		log.Error(err)
	}

	log.Println(*project)

	//m.Model.Save(env)
	m.Model.Save(project)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func DeleteProject(c *gin.Context) {

	id := tools.StringToUint(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  "id format error, req param error,please check",
		})
		return
	}

	// get project name from db
	project := m.NewProject()
	m.Model.First(project, id)
	// check
	if project.Name == "" {
		msg := "db not find id : " + string(id)
		log.Error(msg)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  msg,
		})
		return
	}

	// delete jenkins job
	pj := NewProject()
	pj.AppName = project.Name
	res, msg := pj.DeleteJob()

	if res != "ok" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"res":  "fail",
			"msg":  msg,
		})
		return
	}

	// delete info from table
	m.Model.Delete(project)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
	})

	//project := m.NewProject()
	//if err := c.ShouldBind(project); err != nil {
	//	log.Error(err)
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 0,
	//		"msg":  "fail",
	//	})
	//	return
	//}
	//
	//log.Println("Start calling service operate Jenkins to delete the project: ", project.Name)
	//
	////msg = JenkinsProject(project, "delete")
	////if msg != "ok" {
	////	c.JSON(http.StatusOK, gin.H{
	////		"code": 0,
	////		"msg":  "fail",
	////	})
	////	return
	////}
	//
	//log.Println("Call service operate Jenkins to delete the project: ", project.Name, " ok", "Start to delete data in mysql-table project")
	//m.Model.Delete(project)
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": 0,
	//	"msg":  "ok",
	//})

}

func GetProjectById(c *gin.Context) {

	var project m.Project
	var count int64
	var getByID m.GetByID
	if err := c.ShouldBindUri(&getByID); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code":   0,
			"count":  1,
			"msg":    "failed",
			"status": "failed",
			"data":   project,
		})
		return
	}

	log.Println(getByID.ID)
	m.Model.First(&project, getByID.ID).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"count":  1,
		"msg":    "ok",
		"status": "ok",
		"data":   project,
	})
}

type Project struct {
	AppName    string `json:"app_name"`
	GitAddress string `json:"git_address"`
}

func NewProject() *Project {
	return &Project{}
}

func (p *Project) CreateJob() (res string, msg string) {

	client := resty.New()
	r, e := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"app_name": p.AppName, "git_address": p.GitAddress}).
		Post(config.GetEnv().ServiceCallJenkins)

	if e != nil {
		return "fail", "request service-call-jenkins create job fail"
	}

	if r.StatusCode() != 200 {
		return "fail", "request service-call-jenkins create job fail, rep code is not 200"
	}

	rep := NewCommonRep()
	if e = json.Unmarshal(r.Body(), rep); e != nil {
		return "fail", "json marshal fail"
	}

	return rep.Res, rep.Msg
}

func (p *Project) DeleteJob() (res string, msg string) {
	client := resty.New()
	r, e := client.R().
		SetQueryString("app_name=" + p.AppName).
		Delete(config.GetEnv().ServiceCallJenkins)

	if e != nil {
		return "fail", "request service-call-jenkins delete job fail"
	}

	if r.StatusCode() != 200 {
		return "fail", "request service-call-jenkins delete job fail, rep code is not 200"
	}

	rep := NewCommonRep()
	if e = json.Unmarshal(r.Body(), rep); e != nil {
		return "fail", "json marshal fail"
	}

	return rep.Res, rep.Msg
}

type CommonRep struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Res  string `json:"res"`
}

func NewCommonRep() *CommonRep {
	return &CommonRep{}
}