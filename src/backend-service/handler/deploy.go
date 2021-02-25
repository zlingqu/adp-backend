package handler

import (
	gitlab_svc "app-deploy-platform/3rd-api/gitlab/service"
	"app-deploy-platform/backend-service/config"
	m "app-deploy-platform/backend-service/model"
	"app-deploy-platform/common/tools"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/zuoshenglo/libs/logs/logrus"
	tool "github.com/zuoshenglo/tools"
	"gopkg.in/resty.v1"
)

func DeployOnline(c *gin.Context) {
	var ID m.ID
	if err := c.ShouldBindUri(&ID); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
		})
		return
	}

	// select table deploy
	d := m.NewDeploy()
	m.DB.First(d, ID.ID)
	log.Info("User attempts to deploy : ", d)

	// get env by id
	env, e := getEnvById(d.EnvId)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
		})
		return
	}

	// get project by id
	project, e := getProjectById(d.AppId)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
		})
		return
	}
	project.Data.GitRepository = gitlab_svc.GitlabUrlCheck(project.Data.GitRepository)
	res, msg, url, lb := "", "", "", ""

	log.Info("First request service-operate-jenkins build")
	res, msg, url, lb = ReqServiceOperateJenkinsBuild(env, project, d)

	if res == "fail" {
		// req service-call-jenkins
		log.Info("Failed to request service-operate-jenkins for the first time. Try to request service-call-jenkins trigger.")
		res, msg = ReqServiceCallJenkinsTrigger(project)
		if res == "ok" {
			time.Sleep(2 * time.Second)
			log.Info("The attempt to request service-call-jenkins trigger succeeded, and the second start to request service-operate-jenkins build")
			res, msg, url, lb = ReqServiceOperateJenkinsBuild(env, project, d)
			if res == "fail" {
				msg = "Please confirm whether there is Jenkinsfile under the build branch!"
			}
		}
	}

	if res == "ok" {
		d.LastBuildInfo = lb
		d.Status = "building"
		d.JenkinsBuildToken = url
		t2, _ := time.ParseInLocation("2006-01-02T15:04:05Z", time.Now().Format("2006-01-02T15:04:05Z"), time.Local)
		d.LastDeploy = t2
		m.DB.Save(d)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  res,
		"msg":  msg,
	})
}

func ReqServiceCallJenkinsJobUpdate(appName string, gitAddress string) (res string, msg string) {
	client := resty.New()
	r, e := client.R().SetHeader("Accept", "application/json").
		SetBody(map[string]string{"name": appName, "git_address": gitAddress}).Put(config.GetEnv().ServiceCallJenkinsJobUpdateAddress)
	if e != nil {
		log.Error(e)
		return "fail", "fail"
	}
	if r.StatusCode() != 200 {
		return "fail", "req service-call-jenkins update job code is not 200"
	}

	var serviceCallJenkinsJobUpdateRespone m.ServiceCallJenkinsJobUpdateRespone
	e = json.Unmarshal(r.Body(), &serviceCallJenkinsJobUpdateRespone)
	if e != nil {
		return "fail", "json unmarshal fail"
	}

	return serviceCallJenkinsJobUpdateRespone.Res, serviceCallJenkinsJobUpdateRespone.Msg
}

func ReqServiceCallJenkinsTrigger(p m.GetProjectById) (res string, msg string) {
	client := resty.New()
	r, e := client.R().SetHeader("Accept", "application/json").
		SetQueryString("token=" + p.Data.Name).Post(config.GetEnv().ServiceCallJenkinsTriggerAddress)

	if e != nil {
		log.Error(e)
		return "fail", "fail"
	}

	if r.StatusCode() != 200 {
		return "fail", "req service-call-jenkins response code is not 200"
	}

	var serviceCallJenkinsTriggerRespone m.ServiceCallJenkinsTriggerRespone
	e = json.Unmarshal(r.Body(), &serviceCallJenkinsTriggerRespone)
	if e != nil {
		return "fail", "json unmarshal fail"
	}
	return serviceCallJenkinsTriggerRespone.Res, serviceCallJenkinsTriggerRespone.Msg
}

func ReqServiceOperateJenkinsBuild(env m.GetEnvById, project m.GetProjectById, d *m.Deploy) (res string, msg string, url string, lb string) {
	re, ms, url, lb := "ok", "ok", "", ""
	var reqJenkinsBuild m.ReqJenkinsBuild
	//reqJenkinsBuild.SetReqJenkinsBuildData(env.Data, project.Data, *d).SetReplics(env.Data.Name, project.Data.PodsNum).SetUnityAppName(project.Data.UnityAppId)
	reqJenkinsBuild.SetReqJenkinsBuildData(env.Data, project.Data, *d).SetUnityAppName(project.Data.UnityAppId)
	byte, _ := json.Marshal(reqJenkinsBuild)
	client := resty.New()
	r, e := client.R().SetHeader("Accept", "application/json").SetBody(string(byte)).Post(config.GetEnv().JenkinsBuildAddress)
	if e != nil {
		log.Error(e)
		return "fail", "fail", url, lb
	}

	if r.StatusCode() != 200 {
		return "fail", "Request Jenkins to build, status code is not 200", url, lb
	}

	var jenkinsBuildResponse m.JenkinsBuildResponse
	e = json.Unmarshal(r.Body(), &jenkinsBuildResponse)
	if e != nil {
		log.Error(e)
		return "fail", "json unmarshal fail", url, lb
	}

	if jenkinsBuildResponse.Status == "faild" {
		return "fail", jenkinsBuildResponse.Info, url, lb
	} else {
		url = jenkinsBuildResponse.Url
		lb = string(r.Body())
		return re, ms, url, lb
	}
}

func DeleteDeploy(c *gin.Context) {

	var ID m.ID
	if err := c.ShouldBindUri(&ID); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  "fail",
		})
		return
	}

	d := m.NewDeploy()
	d.ID = ID.ID
	m.DB.Delete(d)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func PostChange(c *gin.Context) {
	var postChange m.PostChange
	if err := c.ShouldBind(&postChange); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"res":  "fail",
		})
		return
	}

	log.Info("postChange : ")
	var d m.Deploy
	m.DB.Model(&d).Where("jenkins_build_token = ?", postChange.Token).Update("status", postChange.Status)
	// 通知的结果保存到db完整后，发送消息通知service-build-status-send 服务。
	res, msg := postServiceBuildStatusSend(postChange.Token, postChange.Status)
	if res != "ok" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"res":  res,
			"msg":  msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "ok",
		"msg":  "ok",
	})
}

// 把构建结果通知前端广播服务，目的让前端知道此服务的结果
func postServiceBuildStatusSend(jenkinsBuildToken string, status string) (string, string) {
	serviceBuildStatusSendUrl := config.GetEnv().ServiceBuildStatusSendUrl
	sendJsonString := fmt.Sprintf(`{"jenkins_build_token": "%s", "status": "%s"}`, jenkinsBuildToken, status)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(sendJsonString).
		Post(serviceBuildStatusSendUrl)

	if err != nil {
		msg := "部署完成，数据存入db，请求service-build-status-send的接口的时候失败！"
		log.Error(msg)
		return "fail", msg
	}

	if resp.StatusCode() != 200 {
		msg := "部署完成，数据存入db，请求service-build-status-send的接口的时候, 返回的状态码不为200。"
		log.Error(msg)
		return "fail", msg
	}

	return "ok", "ok"
}

func PostUpdate(c *gin.Context) {
	var up m.UpdateDeploy

	if err := c.ShouldBind(&up); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"res":  "json转换失败",
		})
		return
	}
	d := m.NewDeploy()
	d.UpdateDeploy = up
	log.Info(up)
	log.Info(d)
	userInfo, e := getOwnerChinaName(up.OwnerEnglishName)
	if e != nil {
		log.Info("find user name error :", e)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  e,
			"res":  e.Error(),
		})
		return
	} else {
		up.OwnerChinaName = userInfo
	}

	// m.DB.Model(d.TableName()).Updates(&up)
	m.DB.Model(d).Updates(*d)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"res":  "ok",
	})
}

func PostDeploy(c *gin.Context) {
	deploy := m.NewDeploy()

	if err := c.ShouldBindJSON(&deploy); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"res":  "fail",
		})
		return
	}

	log.Info("The received front-end request data is:", deploy)
	log.Info("add chinese name")
	userInfo, e := getUserInfo(deploy.OwnerEnglishName)
	if e != nil {
		log.Info("find user name error :", e)
		deploy.OwnerChinaName = ""
	} else {
		deploy.OwnerChinaName = userInfo[deploy.OwnerEnglishName]
	}

	// set default status
	if deploy.UpdateDeploy.EnvId == 11 {
		deploy.Status = "pending"
	} else {
		deploy.Status = "reviewed"
	}

	t2, _ := time.ParseInLocation("2006-01-02T15:04:05Z", time.Now().Format("2006-01-02T15:04:05Z"), time.Local)
	deploy.LastDeploy = t2

	// if deploy.VersionControlMode=="GitCommitId" && deploy.GitCommitId=="last"{
	// 	deploy.GitCommitId= "abc"
	// }

	m.DB.Create(deploy)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "ok",
	})
}

func GetDeploy(c *gin.Context) {

	// init req param
	name := strings.TrimSpace(c.DefaultQuery("name", ""))
	rDeployNamespace := strings.TrimSpace(c.DefaultQuery("deploy_namespace", ""))
	rOwnerName := strings.TrimSpace(c.DefaultQuery("owner_name", ""))
	rAppId := strings.TrimSpace(c.DefaultQuery("app", ""))
	limit, offset := tools.GetMysqlLimitOffset(c.DefaultQuery("page", "1"), c.DefaultQuery("size", "10"))

	var deploy []m.ReqDeploy
	appIdList := make([]int, 0)
	ownerEnglishNameList := make([]string, 0)
	ownerChineseNameList := make([]string, 0)
	var count int64
	db := m.DB
	if rDeployNamespace != "" {
		db = db.Where("k8s_namespace in (?)", strings.Split(rDeployNamespace, ","))
	}

	if rOwnerName != "" {
		db = db.Where("owner_china_name in (?)", strings.Split(rOwnerName, ","))
	}

	if rAppId != "" {
		db = db.Where("app_id in (?)", strings.Split(rAppId, ","))
	}

	//name := c.DefaultQuery("name", "")
	log.Info(fmt.Sprintf("req parms: name: %s, limit: %d, offset: %d", name, limit, offset))

	log.Info("开始请求service-adp-env的数据。")
	reqEnvData, e := getEnv()
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  e,
			"res":  "fail",
		})
		return
	}
	log.Info("请求service-adp-env的数据完成。")

	log.Info("开始请求service-adp-project的数据。")
	reqProjectData, e := getProjectByName(name)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  e,
			"res":  "fail",
		})
		return
	}
	log.Info("请求service-adp-project的数据完成。")

	var reqProjectAllData map[int]string
	if name != "" {
		reqProjectAllData, e = getProjectByName("")
		if e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  e,
				"res":  "fail",
			})
			return
		}
	} else {
		reqProjectAllData = reqProjectData
	}

	log.Info("Request service service-adp-project all data succeeded")

	for k, _ := range reqProjectData {
		//appNameList = append(appNameList, v)
		appIdList = append(appIdList, k)
	}

	if name != "" {
		log.Info("开始请求service-adp-user的数据。")
		requestUserInfo, e := getUserInfoV2(name)
		log.Info("请求service-adp-user的数据完成。")
		if e != nil {
			log.Error(e)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  e,
				"res":  "fail",
			})
			return
		}
		log.Info("Request "+config.GetEnv().SearchUserAddress+" succeeded", requestUserInfo)
		for k, v := range requestUserInfo {
			ownerEnglishNameList = append(ownerEnglishNameList, k)
			ownerChineseNameList = append(ownerChineseNameList, v)
		}

		db = db.Where("name like ? OR app_id in (?) OR owner_english_name in (?) OR owner_china_name in (?)",
			"%"+name+"%", appIdList, ownerEnglishNameList, ownerChineseNameList)

		//db.Where("name like ?", "%" + name + "%").
		//	Or("app_id in (?)", appIdList).
		//	Or("owner_english_name in (?)", ownerEnglishNameList).
		//	Or("owner_china_name in (?)", ownerChineseNameList)
	}

	log.Info("Start querying table deploy")
	log.Info("开始查询数据库的数据。")
	db.Limit(limit).Offset(offset).Find(&deploy)
	log.Info("开始查询数据库的统计数据。")
	db.Model(&m.ReqDeploy{}).Count(&count)

	//log.Println(reqProjectData)
	// add appName, envName, ownerName
	for k, v := range deploy {
		//deploy[k].AppName = strings.Split(reqProjectAllData[int(v.AppId)], "::::::")[0]
		tmpProjectInfolist := strings.Split(reqProjectAllData[int(v.AppId)], "::::::")
		if len(tmpProjectInfolist) == 6 {
			deploy[k].AppName = tmpProjectInfolist[0]
			deploy[k].GitRepository = tmpProjectInfolist[1]
			deploy[k].LanguageType = tmpProjectInfolist[2]
			deploy[k].IfUseModel = tools.StringToBool(tmpProjectInfolist[3])
			deploy[k].IfUseGitManagerModel = tools.StringToBool(tmpProjectInfolist[4])
			deploy[k].ModelGitRepository = tmpProjectInfolist[5]
			deploy[k].EnvName = reqEnvData[int(v.EnvId)]
		}
		//deploy[k].OwnerChinaName = requestUserInfo[v.Owner]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "ok",
		"res":   "ok",
		"count": count,
		"data":  deploy,
	})
}

func PostDeployList(c *gin.Context) {
	var deploy []m.ReqDeploy
	var name m.Name
	//appNameList := make([]string, 0)
	appIdList := make([]int, 0)
	ownerEnglishNameList := make([]string, 0)
	ownerChineseNameList := make([]string, 0)
	var count int64
	//var msg string

	// check band
	if err := c.ShouldBindJSON(&name); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "fail",
			"res":  "fail",
		})
		return
	}

	log.Info("Start requesting environment data")
	reqEnvData, e := getEnv()
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  e,
			"res":  "fail",
		})
		return
	}

	log.Info("Request environment data succeeded")

	log.Info("The parameters requested by the user are: ", name.Name)
	reqProjectData, e := getProjectByName(name.Name)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  e,
			"res":  "fail",
		})
		return
	}

	log.Info("Request service service-adp-project succeeded")

	var reqProjectAllData map[int]string
	if name.Name != "" {
		reqProjectAllData, e = getProjectByName("")
		if e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  e,
				"res":  "fail",
			})
			return
		}
	} else {
		reqProjectAllData = reqProjectData
	}

	log.Info("Request service service-adp-project all data succeeded")

	for k, _ := range reqProjectData {
		//appNameList = append(appNameList, v)
		appIdList = append(appIdList, k)
	}

	if name.Name != "" {
		log.Info(name.Name)
		requestUserInfo, e := getUserInfo(name.Name)
		if e != nil {
			log.Error(e)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  e,
				"res":  "fail",
			})
			return
		}
		log.Info("Request "+config.GetEnv().SearchUserAddress+" succeeded", requestUserInfo)
		for k, v := range requestUserInfo {
			ownerEnglishNameList = append(ownerEnglishNameList, k)
			ownerChineseNameList = append(ownerChineseNameList, v)
		}
		m.DB.Where("name like ?", "%"+name.Name+"%").
			Or("app_id in (?)", appIdList).
			Or("owner_english_name in (?)", ownerEnglishNameList).
			Or("owner_china_name in (?)", ownerChineseNameList).Find(&deploy).Count(&count)
	}

	log.Info("Start querying table deploy")
	if name.Name == "" {
		m.DB.Find(&deploy).Count(&count)
	}

	//log.Println(reqProjectData)
	// add appName, envName, ownerName
	for k, v := range deploy {
		//deploy[k].AppName = strings.Split(reqProjectAllData[int(v.AppId)], "::::::")[0]
		tmpProjectInfolist := strings.Split(reqProjectAllData[int(v.AppId)], "::::::")
		if len(tmpProjectInfolist) == 2 {
			deploy[k].AppName = tmpProjectInfolist[0]
			deploy[k].GitRepository = tmpProjectInfolist[1]
			deploy[k].EnvName = reqEnvData[int(v.EnvId)]
		}
		//deploy[k].OwnerChinaName = requestUserInfo[v.Owner]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "ok",
		"res":   "ok",
		"count": count,
		"data":  deploy,
	})
}

func getUserInfoV2(name string) (map[string]string, error) {
	// init
	userInfoList := make(map[string]string)
	userInfourl := config.GetEnv().ServiceAdpUserUrl + "/api/v1/get-user-for-name"
	var getUserInfo m.GetUserInfo

	// req service-adp-user
	client := resty.New()
	r, e := client.R().SetQueryParam("name", name).SetHeader("Accept", "application/json").Get(userInfourl)

	// dispose
	if e != nil {
		log.Error("请求service-adp-user错误！")
		return userInfoList, e
	}

	if r.StatusCode() != 200 {
		log.Error("请求service-adp-user的返回状态码，不为200！")
		return userInfoList, errors.New("请求service-adp-user的返回状态码，不为200！")
	}

	e = json.Unmarshal(r.Body(), &getUserInfo)
	if e != nil {
		log.Error(e)
		return userInfoList, e
	}

	for _, v := range getUserInfo.Data {
		userInfoList[v["owner_english_name"].(string)] = v["owner_china_name"].(string)
	}

	return userInfoList, nil
}

func getUserInfo(name string) (map[string]string, error) {
	userInfoList := make(map[string]string)
	userInfourl := config.GetEnv().SearchUserAddress
	client := resty.New()
	r, e := client.R().SetQueryParam("name", name).SetHeader("Accept", "application/json").Get(userInfourl)
	if e != nil {
		log.Error(e)
		return userInfoList, e
	}

	if r.StatusCode() != 200 {
		log.Info("request , " + userInfourl + " status code is not 200")
		return userInfoList, errors.New("request , " + userInfourl + " status code is not 200")
	}

	log.Info("The response of mis-admin-backend:")

	var getUserInfo m.GetUserInfo
	e = json.Unmarshal(r.Body(), &getUserInfo)
	if e != nil {
		log.Error(e)
		return userInfoList, e
	}

	for _, v := range getUserInfo.Data {
		userInfoList[v["username"].(string)] = v["name"].(string)
	}

	return userInfoList, nil
}

func getOwnerChinaName(ownerEnglishName string) (string, error) {
	repURL := config.GetEnv().SearchOwnerChinaForEnglishNameUrl

	client := resty.New()
	r, e := client.R().SetQueryParam("ownerEnglishName", ownerEnglishName).Get(repURL)

	if e != nil {
		log.Error(e)
		return "", errors.New("使用用户的英文名请求service-adp-user查询用户的中文名失败错误。")
	}

	if r.StatusCode() != 200 {
		log.Error("请求用户的英文名请求service-adp-user查询用户的中文名，返回的状态码不为200。")
		return "", errors.New("请求用户的英文名请求service-adp-user查询用户的中文名，返回的状态码不为200。")
	}

	var repData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Res  string `json:"res"`
		Data struct {
			OwnerChinaName string `json:"owner_china_name"`
		} `json:"data"`
	}

	e = json.Unmarshal(r.Body(), &repData)

	if e != nil {
		log.Error(e)
		return "", errors.New("请求用户的英文名请求service-adp-user查询用户的中文名, 对返回结果进行反序列化失败。")
	}

	return repData.Data.OwnerChinaName, nil
}

func getEnvById(id uint) (m.GetEnvById, error) {
	var getEnvById m.GetEnvById
	envUrl := config.GetEnv().ServiceEnvAddress
	client := resty.New()
	r, e := client.R().SetHeader("Accept", "application/json").Get(envUrl + "/" + strconv.Itoa(int(id)))
	if e != nil {
		log.Error(e)
		return getEnvById, e
	}

	if r.StatusCode() != 200 {
		log.Info("Request service-adp-env error, status code is not 200")
		return getEnvById, errors.New("request service-adp-env error, status code is not 200")
	}

	e = json.Unmarshal(r.Body(), &getEnvById)
	if e != nil {
		log.Error(e)
		return getEnvById, e
	}

	return getEnvById, nil
}

func getProjectById(id uint) (m.GetProjectById, error) {
	var getProjectById m.GetProjectById
	projectUrl := config.GetEnv().ServiceProjectAddress
	client := resty.New()
	r, e := client.R().SetHeader("Accept", "application/json").Get(projectUrl + "/" + strconv.Itoa(int(id)))
	if e != nil {
		log.Error(e)
		return getProjectById, e
	}

	if r.StatusCode() != 200 {
		log.Info("Request service-adp-project error, status code is not 200")
		return getProjectById, errors.New("request service-adp-project error, status code is not 200")
	}

	e = json.Unmarshal(r.Body(), &getProjectById)
	if e != nil {
		log.Error(e)
		return getProjectById, e
	}

	return getProjectById, nil
}

func getEnv() (map[int]string, error) {
	envList := make(map[int]string, 0)
	envUrl := config.GetEnv().ServiceEnvAddress
	client := resty.New()
	r, e := client.R().SetHeader("Accept", "application/json").Get(envUrl)
	if e != nil {
		log.Error(e)
		return envList, e
	}

	if r.StatusCode() != 200 {
		log.Info("Request service-adp-env error, status code is not 200")
		return envList, errors.New("request service-adp-env error, status code is not 200")
	}

	var postEnv m.PostEnvName
	e = json.Unmarshal(r.Body(), &postEnv)
	if e != nil {
		log.Error(e)
		return envList, e
	}

	for _, v := range postEnv.Data {
		envList[v.ID] = v.Name
	}

	return envList, nil
}

func getProjectByName(name string) (map[int]string, error) {
	projectList := make(map[int]string)
	projectUrl := config.GetEnv().ServiceProjectIdNameGitLangUrl

	client := resty.New()
	r, e := client.R().SetQueryParam("name", name).SetHeader("Accept", "application/json").Get(projectUrl)
	if e != nil {
		log.Error(e)
		return projectList, e
	}

	if r.StatusCode() != 200 {
		log.Info("Request service-adp-project error, status code is not 200")
		return projectList, errors.New("request service-adp-project error, status code is not 200")
	}

	log.Info("The response of service-adp-project is:")

	var getProjectIdNameGitLang m.GetProjectIdNameGitLang
	e = json.Unmarshal(r.Body(), &getProjectIdNameGitLang)
	if e != nil {
		log.Error(e)
		return projectList, e
	}

	for _, v := range getProjectIdNameGitLang.Data {
		projectList[v.ID] = v.Name + "::::::" + v.GitRepository + "::::::" + v.LanguageType + "::::::" + tool.BoolToString(v.IfUseModel) + "::::::" + tool.BoolToString(v.IfUseGitManagerModel) + "::::::" + v.ModelGitRepository
	}

	return projectList, nil
}
