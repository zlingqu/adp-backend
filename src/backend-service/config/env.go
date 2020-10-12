package config

import (
	"app-deploy-platform/backend-service/database"
	"app-deploy-platform/common/tools"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Env struct {
	AppName         string
	Debug           bool
	Database        database.Config
	ServerPort      string
	EventSourcePort string
	RedisIp         string
	RedisPort       string
	RedisPassword   string
	RedisDb         int
	RedisSessionDb  int
	RedisCacheDb    int
	AppSecret       string

	AccessLog     bool
	AccessLogPath string
	ErrorLog      bool
	ErrorLogPath  string
	InfoLog       bool
	InfoLogPath   string

	SqlLog bool

	TemplatePath string // 静态文件相对路径

	JenkinsJobAddress                       string //jenkins构建地址
	JenkinsMultibranchWebhookTriggerAddress string //
	JenkinsBuildAddress                     string

	ServiceProjectAddress              string
	ServiceProjectIdNameGitUrl         string
	ServiceProjectIdNameGitLangUrl     string
	ServiceEnvAddress                  string
	SearchUserAddress                  string
	ServiceCallJenkinsTriggerAddress   string
	ServiceCallJenkinsJobUpdateAddress string
	ServiceAdpUserUrl                  string
	SearchOwnerChinaForEnglishNameUrl  string
	ServiceBuildStatusSendUrl          string
	ServiceCallJenkins                 string
	MisLdapServiceUrl                  string
	ServiceAdpBuildResultUrl           string
	DeployEnv                          string
}

func GetEnv() *Env {
	return &env
}

var (
	localAddr = "http://localhost:" + tools.GetEnvDefault("SERVER_PORT", "80").(string)
	env       = Env{
		AppName: tools.GetEnvDefault("APP_NAME", "service-adp-deploy").(string),
		Debug:   tools.GetEnvDefault("DEBUG_MODEL", false).(bool),

		ServerPort:      tools.GetEnvDefault("SERVER_PORT", "80").(string),
		EventSourcePort: tools.GetEnvDefault("EVENT_SOURCE_PORT", "81").(string),

		Database: database.Config{
			Config: mysql.Config{
				// User:                 "quzl",
				// Passwd:               "quzl",
				// DBName:               "quzl",
				User:                 "adp_test",
				Passwd:               "adp_test",
				DBName:               "test_adp",
				Addr:                 "192.168.3.151:3306",
				Collation:            "utf8mb4_unicode_ci",
				Net:                  "tcp",
				AllowNativePasswords: true,
				ParseTime:            true,
				Loc:                  time.Local,
			},
			MaxOpenConnections: 100,
			MaxIdleConnections: 50,
		},

		RedisIp:       "127.0.0.1",
		RedisPort:     "6379",
		RedisPassword: "",
		RedisDb:       0,

		RedisSessionDb: 1,
		RedisCacheDb:   2,

		AccessLog:     true,
		AccessLogPath: "storage/logs/access.log",

		ErrorLog:     true,
		ErrorLogPath: "storage/logs/error.log",

		InfoLog:     true,
		InfoLogPath: "storage/logs/info.log",

		TemplatePath: "frontend/templates",

		//APP_SECRET: "YbskZqLNT6TEVLUA9HWdnHmZErypNJpL",
		AppSecret:                          "something-very-secret",
		MisLdapServiceUrl:                  tools.GetEnvDefault("MIS_LDAP_SERVICE_URL", "http://mis-ldap-service.mis/search").(string),
		SearchUserAddress:                  tools.GetEnvDefault("SEARCH_USER_ADDRESS", "http://mis-admin-backend.mis.svc.cluster.local/api/open/staff/search").(string),
		ServiceAdpUserUrl:                  tools.GetEnvDefault("SERVER_ADP_USER_URL", localAddr).(string),
		JenkinsJobAddress:                  tools.GetEnvDefault("JENKINS_JOB_ADDRESS", localAddr+"/api/v1/jenkins_job").(string),
		JenkinsBuildAddress:                tools.GetEnvDefault("JENKINS_BUILD_ADDRESS", localAddr+"/api/v1/jenkins/build").(string),
		ServiceProjectIdNameGitUrl:         tools.GetEnvDefault("SERVICE_PROJECT_ID_NAME_GIT", localAddr+"/api/v2/project-id-name-git").(string),
		ServiceProjectIdNameGitLangUrl:     tools.GetEnvDefault("SERVICE_PROJECT_ID_NAME_GIT_LANG_URL", localAddr+"/api/v2/project-id-name-git-lang").(string),
		ServiceProjectAddress:              tools.GetEnvDefault("SERVICE_PROJECT_ADDRESS", localAddr+"/api/v1/project").(string),
		ServiceEnvAddress:                  tools.GetEnvDefault("SERVICE_ENV_ADDRESS", localAddr+"/api/v1/env").(string),
		ServiceCallJenkinsTriggerAddress:   tools.GetEnvDefault("SERVICE_CALL_JENKINS_TRIGGER_ADDRESS", localAddr+"/api/v1/multibranch-webhook-trigger").(string),
		ServiceCallJenkinsJobUpdateAddress: tools.GetEnvDefault("SERVICE_CALL_JENKINS_JOB_UPDATE_ADDRESS", localAddr+"/api/v1/job").(string),
		SearchOwnerChinaForEnglishNameUrl:  tools.GetEnvDefault("SEARCH_OWNER_CHINA_FOR_ENGLISH_NAME_URl", localAddr+"/api/v1/user/get-owner-china-name").(string),
		ServiceBuildStatusSendUrl:          tools.GetEnvDefault("SERVICE_BUILD_STATUS_SEND_URL", "http://service-build-status-send.devops:8080/api/v1/deploy/result").(string),
		ServiceCallJenkins:                 tools.GetEnvDefault("SERVICE_CALL_JENKINS", localAddr+"/api/v1/job").(string),
		ServiceAdpBuildResultUrl:           tools.GetEnvDefault("SERVICE_ADP_BUILD_RESULT", localAddr+"/api/v1/result").(string),
		DeployEnv:                          strings.ToLower(tools.GetEnvDefault("DEPLOY_ENV", "prd").(string)),
	}
)
