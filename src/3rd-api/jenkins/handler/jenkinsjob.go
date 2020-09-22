package handler

import (
	"app-deploy-platform/3rd-api/jenkins/config"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	glibs "github.com/zuoshenglo/go-base-libs"
	"net/http"
)

type OperateJenkins struct {
	DmaiBaseDevopsHttp
	JenkinsJob
}

// 未发现调用接口 调用接口代码已注释
func NewOperateJenkins(c *gin.Context) {
	returnDate := gin.H{
		"status": "faild",
		"info":   "init",
	}

	userJson := map[string]interface{}{}

	if err := c.ShouldBind(&userJson); err != nil {
		log.Error(err)
		returnDate["info"] = "Request parameter error"
		c.JSON(http.StatusOK, returnDate)
		return
	}

	//cfgFile, err := ioutil.ReadFile(conf.ServiceConf.JenkinsJobCfgFile)
	//if err != nil {
	//	log.Error(err)
	//}

	// send to jenkins
	appName := userJson["app_name"].(string)
	action := userJson["action"].(string)
	jenkinsCreateJobUrl := config.ServiceConf.Jenkins.Url + "/createItem?name=" + appName
	jenkinsDeleteJobUrl := config.ServiceConf.Jenkins.Url + "/job/" + appName + "/doDelete"
	jenkinsCfgFile := JenkinsJobCfgFile(appName, userJson["git_address"].(string))

	log.Info("appName: ", appName)
	log.Info("gitAddress: ", userJson["git_address"].(string))
	log.Info("jenkinsFileContent: ", jenkinsCfgFile)

	var req *glibs.HttpRequestCustom
	if action == "create" {
		log.Info("jenkins create job url:", jenkinsCreateJobUrl)
		req = glibs.NewHttpRequestCustom([]byte(jenkinsCfgFile), "POST", jenkinsCreateJobUrl).SetRequestProtocol("http").SetContentType("text/xml")
	}

	if action == "delete" {
		log.Info("jenkins delete job url:", jenkinsDeleteJobUrl)
		req = glibs.NewHttpRequestCustom([]byte(""), "POST", jenkinsDeleteJobUrl).SetRequestProtocol("http").SetContentType("")
	}

	req.SetBasicAuth(config.ServiceConf.Jenkins.User, config.ServiceConf.Jenkins.Password)
	result, err := req.ExecRequest()
	if err != nil {
		log.Error(err)
	}

	//
	if action == "create" {
		if result != "" {
			returnDate["info"] = "project exists"
			c.JSON(http.StatusOK, returnDate)
			return
		} else {
			returnDate["status"] = "ok"
			returnDate["info"] = "create project ok"
		}
	}

	// index jenkins project
	if action == "create" {
		log.Info("re-index project")
		reIndexUrl := config.ServiceConf.Jenkins.Url + "/job/" + appName + "/build"
		req = glibs.NewHttpRequestCustom([]byte(""), "POST", reIndexUrl).SetRequestProtocol("http").SetContentType("text/xml")
		req.SetBasicAuth(config.ServiceConf.Jenkins.User, config.ServiceConf.Jenkins.Password)
		ires, ierr := req.ExecRequest()

		if ierr != nil {
			log.Error(ierr)
			returnDate["status"] = "faild"
			returnDate["info"] = fmt.Sprintf("%s", ierr)
		}

		log.Info(ires)
		//if ires != "" {
		//	log.Error("re-index project faild!!")
		//	log.Info(ires)
		//	returnDate["status"] = "faild"
		//	returnDate["info"] = "re-index project faild!!"
		//	}
	}

	if action == "delete" {
		returnDate["status"] = "ok"
		returnDate["info"] = "delete project ok"
	}

	c.JSON(http.StatusOK, returnDate)
}
