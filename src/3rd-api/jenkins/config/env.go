package config

import (
	"app-deploy-platform/common/tools"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Env struct {
	JenkinsAddress                          string
	JenkinsJobUpdateAddress                 string
	JenkinsMultibranchWebhookTriggerAddress string
	Debug                                   bool
	ServerPort                              string
}

func GetEnv() *Env {
	return &env
}

// 本文件建议在代码协同工具(git/svn等)中忽略

var env = Env{
	Debug:                                   tools.GetEnvDefault("DEBUG_MODEL", true).(bool),
	ServerPort:                              tools.GetEnvDefault("SERVER_PORT", "80").(string),
	JenkinsAddress:                          tools.GetEnvDefault("JENKINS_ADDRESS", "http://jenkins.ops.dm-ai.cn").(string),
	JenkinsJobUpdateAddress:                 tools.GetEnvDefault("JENKINS_JOB_UPDATE_ADDRESS", "http://jenkins.ops.dm-ai.cn/job/").(string),
	JenkinsMultibranchWebhookTriggerAddress: tools.GetEnvDefault("JENKINS_JOB_ADDRESS", "http://jenkins.ops.dm-ai.cn/multibranch-webhook-trigger/invoke").(string),
}

var ServiceConf conf

func init() {
	ServiceConf.init()
}

func getCwd() string {
	dir, _ := os.Getwd()
	return dir
}

type conf struct {
	SeelogCfgFiles          string
	KubernetesCfgFile       string
	MasterKubernetesCfgFile string
	DevKubernetesCfgFile    string
	TestKubernetesCfgFile   string
	StageKubernetesCfgFile  string
	JenkinsJobCfgFile       string
	Jenkins                 jenkins `yaml: "jenkins, omitempty"`
}

type jenkins struct {
	Url             string `yaml: "url, omitempty"`
	User            string `yaml: "user, omitempty"`
	Password        string `yaml: "password, omitempty"`
	Pipelinebaseurl string `yaml: "pipelinebaseurl, omitempty"`
}

func (c *conf) init() {
	c.getConf()
	c.getSeelogCfgFiles()
	c.getKubernetesCfgFile()
}

func (c *conf) getSeelogCfgFiles() *conf {
	c.SeelogCfgFiles = getCwd() + "/conf/seelog.xml"
	log.Println("read seelog config success")
	return c
}

func (c *conf) getKubernetesCfgFile() *conf {
	configPath := getCwd() + "/conf/"
	c.KubernetesCfgFile = configPath + "config"
	c.MasterKubernetesCfgFile = configPath + "master"
	c.DevKubernetesCfgFile = configPath + "dev"
	c.TestKubernetesCfgFile = configPath + "test"
	c.StageKubernetesCfgFile = configPath + "stage"
	c.JenkinsJobCfgFile = configPath + "config.xml"
	log.Println("get k8s config success")
	return c
}

func (c *conf) GetK8sCfg(k8sEnvName string) string {
	switch k8sEnvName {
	case "master":
		return c.MasterKubernetesCfgFile
	case "dev":
		return c.DevKubernetesCfgFile
	case "test":
		return c.TestKubernetesCfgFile
	case "stage":
		return c.StageKubernetesCfgFile
	default:
		return c.DevKubernetesCfgFile
	}
}

func (c *conf) getConf() *conf {
	log.Println("读取service的配置信息。")
	dir, _ := os.Getwd()

	ymalFile, err := ioutil.ReadFile(dir + "/conf/service.yml")
	if err != nil {
		log.Println(fmt.Sprintf("读取项目的配置文件->%s->失败", ymalFile), err)
		panic(err)
	}

	err = yaml.Unmarshal(ymalFile, c)
	if err != nil {
		log.Println("反序列化servcer的配置文件失败", err)
		panic(err)
	}
	log.Println("read service config success")
	return c
}
