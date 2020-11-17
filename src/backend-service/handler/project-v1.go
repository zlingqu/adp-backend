package handler

import (
	"app-deploy-platform/backend-service/config"
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/common/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type Project struct {
	AppName    string `json:"app_name"`
	GitAddress string `json:"git_address"`
}

func NewProject() *Project {
	return &Project{}
}

func (p *Project) CreateJob() (res string, msg string) {

	client := resty.New()
	//调用jenkins接口，创建项目
	r, e := client.R().SetHeader("Content-Type", "application/json").SetBody(map[string]string{"app_name": p.AppName, "git_address": p.GitAddress}).Post(config.GetEnv().ServiceCallJenkins)

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
	r, e := client.R().SetQueryString("app_name=" + p.AppName).Delete(config.GetEnv().ServiceCallJenkins)

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

func PostProject(c *gin.Context) {
	project := m.NewProject()
	var msg string
	var res string

	//打印body
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	fmt.Printf("请求的body原始格式是：%s", string(bodyBytes))
	//打印body

	if err := c.ShouldBindJSON(&project); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{

			"code": 0,
			"msg":  "fail",
		})
		return
	}

	// log.Println("The received front-end request data is : ", project, " , Start to request service-call-jenkins to create the project : ", project.Name)

	pj := NewProject()
	pj.AppName = project.Name
	if strings.HasSuffix(project.GitRepository, ".git") { //如果url写错，没有以.git结尾，将其加上
		project.GitRepository = project.GitRepository + ".git"
	}
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

}
