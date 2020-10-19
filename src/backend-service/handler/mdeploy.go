package handler

import (
	m "app-deploy-platform/backend-service/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MdeployStruct struct {
	Ids   []uint `json:"ids" binding:"required"`
	EnvID uint   `json:"envid" binding:"required"`
	Tag   string `json:"tag" binding:"required"`
}

func Mdeploy(c *gin.Context) {
	MdeployJson := &MdeployStruct{}
	if err := c.BindJSON(&MdeployJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"count": 0,
			"code":  400,
			"res":   "error",
			"msg":   "请求体格式错误，请检查",
		})
		return
	}

	for i := 0; i < len(MdeployJson.Ids); i++ {

		deployID := MdeployJson.Ids[i]

		er := CallJenkinsApiByID(MdeployJson.EnvID, deployID, MdeployJson.Tag)
		if er == "error" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"res":  "fail",
				"msg":  "失败",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"count": len(MdeployJson.Ids),
		"code":  0,
		"res":   "ok",
		"msg":   "构建成功",
	})

}

func CallJenkinsApiByID(envID, id uint, tag string) string {

	// select table deploy
	d := m.NewDeploy()
	m.Model.First(d, id)
	log.Info("User attempts to deploy : ", d)

	// get env by id
	env, e := getEnvById(envID)
	if e != nil {
		return "error"
	}

	// get project by id
	project, e := getProjectById(d.AppId)
	if e != nil {
		return "error"
	}

	if tag != "" {
		d.VersionControlMode = "GitTags"
		d.GitTag = tag
	}

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
		m.Model.Save(d)
		return "ok"
	}
	fmt.Println(msg)
	return "error"
}

// func ReqServiceOperateJenkinsBuild(env m.GetEnvById, project m.GetProjectById, d *m.Deploy) (res string, msg string, url string, lb string) {
// 	re, ms, url, lb := "ok", "ok", "", ""
// 	var reqJenkinsBuild m.ReqJenkinsBuild
// 	reqJenkinsBuild.SetReqJenkinsBuildData(env.Data, project.Data, *d).SetReplics(env.Data.Name, project.Data.PodsNum).SetUnityAppName(project.Data.UnityAppId)
// 	byte, _ := json.Marshal(reqJenkinsBuild)
// 	client := resty.New()
// 	r, e := client.R().SetHeader("Accept", "application/json").
// 		SetBody(string(byte)).Post(config.GetEnv().JenkinsBuildAddress)
// 	if e != nil {
// 		log.Error(e)
// 		return "fail", "fail", url, lb
// 	}

// 	if r.StatusCode() != 200 {
// 		return "fail", "Request Jenkins to build, status code is not 200", url, lb
// 	}

// 	var jenkinsBuildResponse m.JenkinsBuildResponse
// 	e = json.Unmarshal(r.Body(), &jenkinsBuildResponse)
// 	if e != nil {
// 		log.Error(e)
// 		return "fail", "json unmarshal fail", url, lb
// 	}

// 	if jenkinsBuildResponse.Status == "faild" {
// 		return "fail", jenkinsBuildResponse.Info, url, lb
// 	} else {
// 		url = jenkinsBuildResponse.Url
// 		lb = string(r.Body())
// 		return re, ms, url, lb
// 	}
// }
