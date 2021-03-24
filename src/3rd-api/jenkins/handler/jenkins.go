package handler

import (
	"app-deploy-platform/3rd-api/jenkins/config"
	m "app-deploy-platform/3rd-api/jenkins/model"
	svc "app-deploy-platform/3rd-api/jenkins/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

func PostJob(c *gin.Context) {

	jenkins := NewJenkins().SetJenkinsUrl(config.GetEnv().JenkinsAddress)
	if err := c.ShouldBind(jenkins); err != nil || jenkins.AppName == "" || jenkins.GitAddress == "" {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Request parameter error",
			"res":  "fail",
		})
		return
	}

	log.Println("Jenkins create param : ", jenkins)

	// create job
	res, msg := jenkins.CreateJenkinsJob()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
		"res":  res,
	})
}

func DeleteJob(c *gin.Context) {
	appName := c.DefaultQuery("app_name", "app_name")
	if appName == "app_name" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Request parameter error",
			"res":  "fail",
		})
		return
	}

	jenkins := NewJenkins().SetJenkinsUrl(config.GetEnv().JenkinsAddress)
	jenkins.AppName = appName

	res, msg := jenkins.DeleteJenkinsJob()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
		"res":  res,
	})
}

func PutJob(c *gin.Context) {
	var app m.App
	if err := c.ShouldBind(&app); err != nil || app.Name == "" || app.GitAddress == "" {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Request parameter error",
			"res":  "fail",
		})
		return
	}

	res, msg := svc.UpdateJenkinsJobConfig(app.Name, app.GitAddress)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  res,
		"res":  msg,
	})

}

func PostMultibranchWebhookTrigger(c *gin.Context) {
	appName := c.DefaultQuery("token", "defaultApp")
	if appName == "defaultApp" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Request parameter error",
			"res":  "fail",
		})
		return
	}

	//
	log.Println("begin call jenkins")

	client := resty.New()
	r, e := client.R().SetQueryString("token="+appName).SetHeader("Accept", "application/json").Post(config.GetEnv().JenkinsMultibranchWebhookTriggerAddress)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Request jenkins error",
			"res":  "fail",
		})
		return
	}

	if r.StatusCode() != 200 {
		log.Println("Request Jenkins, status code is not 200")
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Request jenkins error, status code is not 200",
			"res":  "fail",
		})
		return
	}

	log.Println(r)
	var multibranchWebhookTrigger m.MultibranchWebhookTrigger
	e = json.Unmarshal(r.Body(), &multibranchWebhookTrigger)

	if e != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "json Unmarshal fail",
			"res":  "fail",
		})
		return
	}

	triggerResults := multibranchWebhookTrigger.Data.TriggerResults
	log.Println(triggerResults)
	log.Println(triggerResults["ANY"])
	log.Println(triggerResults[appName])
	if triggerResults["ANY"] != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  triggerResults["ANY"].(string),
			"res":  "fail",
		})
		return
	} else {
		//td := triggerResults[appName].(m.TriggerResultsData)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
			"res":  "ok",
			"data": triggerResults[appName],
		})
	}
}

type MultibranchWebhookTrigger struct {
	Status string `json:"status"`
	Data   struct {
		TriggerResults map[string]interface{}
	} `json:"data"`
}

type Jenkins struct {
	AppName    string `json:"app_name"`
	GitAddress string `json:"git_address"`
	Url        string `json:"url"` // jenkins url
}

// init
func NewJenkins() *Jenkins {
	return &Jenkins{}
}

func (j *Jenkins) SetJenkinsUrl(url string) *Jenkins {
	j.Url = url
	return j
}

func (j *Jenkins) CreateJenkinsJob() (res string, msg string) {
	createUrl := j.createJobUrl()
	jenkinsCfgFile := j.jenkinsJobCfgFile(j.AppName, j.GitAddress)
	client := resty.New()
	r, e := client.R().SetQueryString("name="+j.AppName).SetHeader("Content-Type", "text/xml").SetBody([]byte(jenkinsCfgFile)).Post(createUrl)

	if e != nil {
		log.Error(e)
		return "fail", "Failed to request Jenkins to create job"
	}

	if r.StatusCode() != 200 {
		return "fail", "Failed to request Jenkins to create job, status code is not 200, app already exists?"
	}

	if string(r.Body()) != "" {
		return "fail", "Failed to request Jenkins to create job, body is not null"
	}

	re, ms := j.MultibranchWebhookTrigger()
	if re != "ok" {
		return re, ms
	}
	return re, ms
}

func (j *Jenkins) DeleteJenkinsJob() (res string, msg string) {
	deleteUrl := j.deleteJobUrl()

	log.Println(deleteUrl)
	client := resty.New()
	r, _ := client.R().Post(deleteUrl)

	if string(r.Body()) != "" {
		return "fail", "delete jenkins job : " + j.AppName + " fail "
	} else {
		return "ok", "ok"
	}
}

func (j *Jenkins) MultibranchWebhookTrigger() (res string, msg string) {

	multibranchWebhookTriggerUrl := j.multibranchWebhookTriggerUrl()
	log.Println("Begin trigger jenkins app : ", j.AppName)
	client := resty.New()
	r, e := client.R().SetQueryString("token="+j.AppName).SetHeader("Accept", "application/json").Post(multibranchWebhookTriggerUrl)

	if e != nil {
		log.Error(e)
		return "fail", "trigger jenkins app : " + j.AppName
	}

	if r.StatusCode() != 200 {
		return "fail", "Failed to trigger jenkins app : " + j.AppName + ", status code is not 200"
	}

	log.Println("trigger jenkins app : " + j.AppName + ", rep : " + string(r.Body()))

	var multibranchWebhookTrigger MultibranchWebhookTrigger
	if e = json.Unmarshal(r.Body(), &multibranchWebhookTrigger); e != nil {
		return "fail", "marshal json fail"
	}

	triggerResults := multibranchWebhookTrigger.Data.TriggerResults
	if triggerResults["ANY"] != nil {
		return "fail", triggerResults["ANY"].(string)
	}

	return "ok", "ok"
}

func (j *Jenkins) multibranchWebhookTriggerUrl() string {
	return j.Url + "/multibranch-webhook-trigger/invoke"
}

func (j *Jenkins) createJobUrl() string {
	return j.Url + "/createItem"
}

func (j *Jenkins) deleteJobUrl() string {
	return j.Url + "/job/" + j.AppName + "/doDelete"
}

func (j *Jenkins) jenkinsJobCfgFile(appName string, gitAddress string) string {

	return fmt.Sprintf(`<?xml version='1.1' encoding='UTF-8'?>
<org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject plugin="workflow-multibranch@2.21">
  <actions/>
  <description></description>
  <properties>
    <org.csanchez.jenkins.plugins.kubernetes.KubernetesFolderProperty plugin="kubernetes@1.15.4">
      <permittedClouds/>
    </org.csanchez.jenkins.plugins.kubernetes.KubernetesFolderProperty>
    <org.jenkinsci.plugins.pipeline.modeldefinition.config.FolderConfig plugin="pipeline-model-definition@1.3.8">
      <dockerLabel></dockerLabel>
      <registry plugin="docker-commons@1.14"/>
    </org.jenkinsci.plugins.pipeline.modeldefinition.config.FolderConfig>
  </properties>
  <folderViews class="jenkins.branch.MultiBranchProjectViewHolder" plugin="branch-api@2.4.0">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </folderViews>
  <healthMetrics>
    <com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric plugin="cloudbees-folder@6.8">
      <nonRecursive>false</nonRecursive>
    </com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
  </healthMetrics>
  <icon class="jenkins.branch.MetadataActionFolderIcon" plugin="branch-api@2.4.0">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </icon>
  <orphanedItemStrategy class="com.cloudbees.hudson.plugins.folder.computed.DefaultOrphanedItemStrategy" plugin="cloudbees-folder@6.8">
    <pruneDeadBranches>true</pruneDeadBranches>
    <daysToKeep>7</daysToKeep>
    <numToKeep>10</numToKeep>
  </orphanedItemStrategy>
  <triggers>
    <com.cloudbees.hudson.plugins.folder.computed.PeriodicFolderTrigger plugin="cloudbees-folder@6.8">
      <spec>H/15 * * * *</spec>
      <interval>3600000</interval>
    </com.cloudbees.hudson.plugins.folder.computed.PeriodicFolderTrigger>
    <com.igalg.jenkins.plugins.mswt.trigger.ComputedFolderWebHookTrigger plugin="multibranch-scan-webhook-trigger@1.0.1">
      <spec></spec>
      <token>%s</token>
    </com.igalg.jenkins.plugins.mswt.trigger.ComputedFolderWebHookTrigger>
  </triggers>
  <disabled>false</disabled>
  <sources class="jenkins.branch.MultiBranchProject$BranchSourceList" plugin="branch-api@2.4.0">
    <data>
      <jenkins.branch.BranchSource>
        <source class="jenkins.plugins.git.GitSCMSource" plugin="git@3.10.0">
          <id>%s-%d</id>
          <remote>%s</remote>
          <credentialsId>devops-use</credentialsId>
		  <traits>
		    <jenkins.plugins.git.traits.BranchDiscoveryTrait/>
			<jenkins.plugins.git.traits.SubmoduleOptionTrait>
				<extension class="hudson.plugins.git.extensions.impl.SubmoduleOption">
					<disableSubmodules>false</disableSubmodules>
					<recursiveSubmodules>true</recursiveSubmodules>
					<trackingSubmodules>false</trackingSubmodules>
					<reference></reference>
					<parentCredentials>true</parentCredentials>
				</extension>
			</jenkins.plugins.git.traits.SubmoduleOptionTrait>
          </traits>
        </source>
        <strategy class="jenkins.branch.DefaultBranchPropertyStrategy">
          <properties class="empty-list"/>
        </strategy>
      </jenkins.branch.BranchSource>
    </data>
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </sources>
  <factory class="org.jenkinsci.plugins.workflow.multibranch.WorkflowBranchProjectFactory">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
    <scriptPath>Jenkinsfile</scriptPath>
  </factory>
</org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject>
`, appName, appName, time.Now().UnixNano(), gitAddress)
}
