package service

import (
	"app-deploy-platform/3rd-api/jenkins/config"
	m "app-deploy-platform/3rd-api/jenkins/model"

	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

func UpdateJenkinsJobConfig(appName string, appGitAddress string) (res string, msg string) {
	jenkinsJobUpdateUrl := config.GetEnv().JenkinsJobUpdateAddress + appName + "/config.xml"
	jenkinsCfgFile := m.JenkinsJobCfgFile(appName, appGitAddress)
	client := resty.New()
	r, e := client.R().SetHeader("Content-Type", "text/xml").SetBody([]byte(jenkinsCfgFile)).Post(jenkinsJobUpdateUrl)

	if e != nil {
		log.Error(e)
		return "fail", "Failed to request Jenkins to update configuration"
	}

	if r.StatusCode() != 200 {
		return "fail", "Failed to request Jenkins to update configuration, status code is not 200"
	}

	if string(r.Body()) != "" {
		return "fail", "Failed to request Jenkins to update configuration, body is not null"
	}

	return "ok", "ok"
}
