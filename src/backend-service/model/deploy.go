package model

import (
	"app-deploy-platform/common/tools"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/zuoshenglo/libs/logs/logrus"
)

type UpdateDeploy struct {
	AppId              uint   `json:"app_id"`
	Branch             string `json:"branch" gorm:"type:varchar(80)"`
	EnvId              uint   `json:"env_id"`
	GitCommitId        string `json:"git_commit_id" gorm:"type:varchar(150)"`
	GitTag             string `json:"git_tag" gorm:"type:varchar(150)"`
	ID                 uint   `json:"id" gorm:"primary_key"`
	Name               string `json:"name" gorm:"type:varchar(80)"`
	OwnerEnglishName   string `json:"owner_english_name" gorm:"type:varchar(80)"`
	OwnerChinaName     string `json:"owner_china_name" gorm:"type:varchar(80)"`
	Status             string `json:"status" gorm:"type:varchar(20)"`
	VersionControlMode string `json:"version_control_mode" gorm:"type:varchar(80)"`
	ApolloClusterName  string `json:"apollo_cluster_name" gorm:"type:varchar(80)"`
	ApolloNamespace    string `json:"apollo_namespace" gorm:"type:varchar(80)"`
	K8sNamespace       string `json:"k8s_namespace" gorm:"type:varchar(80)"`
	JsVersion          string `json:"js_version" gorm:"type:varchar(80)"`
	ModelBranch        string `json:"model_branch" gorm:"type:varchar(200)"`
}

type Deploy struct {
	UpdateDeploy
	//ID                  uint      `json:"id" gorm:"primary_key"`
	CreatedAt  time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastDeploy time.Time `json:"last_deploy"`
	//Name                string    `json:"name" gorm:"type:varchar(80);unique_index"`
	//Name                string    `json:"name" gorm:"type:varchar(80)"`
	App int `json:"app"`
	//AppId               uint       `json:"app_id"`
	Owner string `json:"owner" gorm:"type:varchar(80)"`
	//OwnerEnglishName    string    `json:"owner_english_name" gorm:"type:varchar(80)"`
	//OwnerChinaName      string    `json:"owner_china_name" gorm:"type:varchar(80)"`
	//Branch              string    `json:"branch" gorm:"type:varchar(80)"`
	Env int `json:"env"`
	//EnvId               uint       `json:"env_id"`
	Version string `json:"version" gorm:"type:varchar(80)"`
	//VersionControlMode  string    `json:"version_control_mode" gorm:"type:varchar(80)"`
	//GitCommitId         string    `json:"git_commit_id" gorm:"type:varchar(150)"`
	//GitTag              string    `json:"git_tag" gorm:"type:varchar(150)"`
	//Status              string    `json:"status" gorm:"type:varchar(20)"`
	LastBuildInfo     string `json:"last_build_info" gorm:"type:text"`
	JenkinsBuildToken string `json:"jenkins_build_token" gorm:"type:text"`
}

type ReqDeploy struct {
	Deploy
	AppName              string `json:"app_name"`
	EnvName              string `json:"env_name"`
	GitRepository        string `json:"git_repository"`
	LanguageType         string `json:"language_type"`
	IfUseModel           bool   `json:"if_use_model"`
	IfUseGitManagerModel bool   `json:"if_use_git_manager_model"`
	ModelGitRepository   string `json:"model_git_repository"`
}

func (Deploy) TableName() string {
	return "deploy"
}

func NewDeploy() *Deploy {
	d := &Deploy{}
	// if !Model.HasTable(d.TableName()) {
	// 	Model.CreateTable(d)
	// }
	if Model.HasTable(d.TableName()) { //判断表是否存在
		Model.AutoMigrate(d) //存在就自动适配表，也就说原先没字段的就增加字段
	} else {
		Model.CreateTable(d) //不存在就创建新表
	}
	return d
}

// post list
type Name struct {
	Name string `json:"name"`
}

type PostProjectName struct {
	Code  int                `json:"code"`
	Count int                `json:"count"`
	Data  []ReqProjectResult `json:"data"`
	Msg   string             `json:"msg"`
}

type GetProjectIdNameGit struct {
	Code int                `json:"code"`
	Data []ReqProjectResult `json:"data"`
	Res  string             `json:"res"`
	Msg  string             `json:"msg"`
}

type GetProjectIdNameGitLang struct {
	Code int                `json:"code"`
	Data []ReqProjectResult `json:"data"`
	Res  string             `json:"res"`
	Msg  string             `json:"msg"`
}

type GetProjectById struct {
	Code  int     `json:"code"`
	Count int     `json:"count"`
	Data  Project `json:"data"`
	Msg   string  `json:"msg"`
}

type PostEnvName struct {
	Code  int            `json:"code"`
	Count int            `json:"count"`
	Data  []ReqEnvResult `json:"data"`
	Msg   string         `json:"msg"`
}

type GetEnvById struct {
	Code  int    `json:"code"`
	Count int    `json:"count"`
	Data  Env    `json:"data"`
	Msg   string `json:"msg"`
}

type ReqProjectResult struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GitRepository        string `json:"git_repository"`
	LanguageType         string `json:"language_type"`
	IfUseModel           bool   `json:"if_use_model"`
	IfUseGitManagerModel bool   `json:"if_use_git_manager_model"`
	ModelGitRepository   string `json:"model_git_repository"`
}

type ReqEnvResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetUserInfo struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data []map[string]interface{} `json:"data"`
}

type ID struct {
	ID uint `uri:"id" binding:"required"`
}

type ReqJenkinsBuild struct {
	GitAddress                string `json:"git_address"`
	AppName                   string `json:"app_name"`
	ProductName               string `json:"product_name"`
	CodeLanguage              string `json:"code_language"`
	IfAddUnityProject         bool   `json:"if_add_unity_project"`
	UnityAppName              string `json:"unity_app_name"`
	DeployEnvType             string `json:"deploy_env_type"`
	IfCompile                 bool   `json:"if_compile"`
	IfCompileCache            bool   `json:"if_compile_cache"`
	IfCompileParam            bool   `json:"if_compile_param"`
	CompileParam              string `json:"compile_param"`
	IfCompileImage            bool   `json:"if_compile_image"`
	CompileImage              string `json:"compile_image"`
	IfMakeImage               bool   `json:"if_make_image"`
	IfUseDomainName           bool   `json:"if_use_domain_name"`
	DomainName                string `json:"domain_name"`
	IfUseHttps                bool   `json:"if_use_https"`
	IfUseHttp                 bool   `json:"if_use_http"`
	IfUseGrpc                 bool   `json:"if_use_grpc"`
	IfUseGbs                  bool   `json:"if_use_gbs"`
	IfUseSticky               bool   `json:"if_use_sticky"`
	IfDeploy                  bool   `json:"if_deploy"`
	IfUseModel                bool   `json:"if_use_model"`
	IfUseGitManagerModel      bool   `json:"if_use_git_manager_model"`
	ModelGitRepository        string `json:"model_git_repository"`
	IfSaveModelBuildComputer  bool   `json:"if_save_model_build_computer"`
	IfUseConfigmap            bool   `json:"if_use_configmap"`
	IfUseAutoDeployFile       bool   `json:"if_use_auto_deploy_file"`
	AutoDeployContent         string `json:"auto_deploy_content"`
	IfUseCustomDockerfile     bool   `json:"if_use_custom_dockerfile"`
	IfUseRootDockerfile       bool   `json:"if_use_root_dockerfile"`
	DockerfileContent         string `json:"dockerfile_content"`
	ServeType                 string `json:"serve_type"`
	ReplicationControllerType string `json:"replication_controller_type"`
	IfUseGpuCard              bool   `json:"if_use_gpu_card"`
	GpuControlMode            string `json:"gpu_control_mode"`
	GpuCardCount              int    `json:"gpu_card_count"`
	GpuMemCount               int    `json:"gpu_mem_count"`
	BranchName                string `json:"branch_name"`
	Version                   string `json:"version"`
	VersionControlMode        string `json:"version_control_mode"`
	GitCommitId               string `json:"git_commit_id"`
	GitTag                    string `json:"git_tag"`
	ApolloClusterName         string `json:"apollo_cluster_name"`
	ApolloNamespace           string `json:"apollo_namespace"`
	DeployEnv                 string `json:"deploy_env"`
	DeployEnvStatus           string `json:"deploy_env_status"`
	Replics                   int    `json:"replics"`
	ContainerPort             int    `json:"container_port"`
	ServiceListenPort         string `json:"service_listen_port"`
	CpuRequest                string `json:"cpu_request"`
	CpuLimit                  string `json:"cpu_limit"`
	MemoryRequest             string `json:"memory_request"`
	MemoryLimit               string `json:"memory_limit"`
	IfStorageLocale           bool   `json:"if_storage_locale"`
	StoragePath               string `json:"storage_path"`
	IfCheckPodsStatus         bool   `json:"if_check_pods_status"`
	IfUseIstio                bool   `json:"if_use_istio"`
	IfUseApolloOfflineEnv     bool   `json:"if_use_apollo_offline_env"`
	JsVersion                 string `json:"js_version"`
	ModelBranch               string `json:"model_branch"`
}

func (r *ReqJenkinsBuild) SetReqJenkinsBuildData(env Env, project Project, d Deploy) *ReqJenkinsBuild {
	r.GitAddress = project.GitRepository
	r.AppName = project.Name
	r.ProductName = d.K8sNamespace
	r.CodeLanguage = project.LanguageType
	r.IfAddUnityProject = *project.IfAddUnityProject
	r.DeployEnvType = project.DeployEnvType
	r.IfCompile = *project.IfCompile
	r.IfCompileCache = *project.IfCompileCache
	r.IfCompileParam = *project.IfCompileParam
	r.CompileParam = project.CompileParam
	r.IfCompileImage = *project.IfCompileImage
	r.CompileImage = project.CompileImage
	r.IfMakeImage = *project.IfMakeImage
	r.IfUseDomainName = *project.IfUseDomainName
	r.DomainName = project.DomainName
	r.IfUseHttps = *project.IfUseHttps
	r.IfUseGrpc = *project.IfUseGrpc
	r.IfUseGbs = *project.IfUseGbs
	r.IfUseSticky = *project.IfUseSticky
	r.IfUseHttp = *project.IfUseHttp
	r.IfDeploy = *project.IfDeploy
	r.IfUseModel = *project.IfUseModel
	r.IfUseGitManagerModel = *project.IfUseGitManagerModel
	r.ModelGitRepository = project.ModelGitRepository
	r.IfSaveModelBuildComputer = *project.IfSaveModelBuildComputer
	r.IfUseConfigmap = *project.IfUseConfigmap
	r.IfUseAutoDeployFile = *project.IfUseAutoDeployFile
	r.AutoDeployContent = project.AutoDeployContent
	r.IfUseCustomDockerfile = *project.IfUseCustomDockerfile
	r.IfUseRootDockerfile = *project.IfUseRootDockerfile
	r.DockerfileContent = project.DockerfileContent
	r.ServeType = project.ServeType
	r.ReplicationControllerType = project.ReplicationControllerType
	r.IfUseGpuCard = *project.IfUseGpuCard
	r.GpuControlMode = project.GpuControlMode
	r.GpuCardCount = project.GpuCardCount
	r.GpuMemCount = project.GpuMemCount
	r.BranchName = d.Branch
	r.Version = d.Version
	r.VersionControlMode = d.VersionControlMode
	r.GitCommitId = d.GitCommitId
	r.GitTag = d.GitTag
	r.ApolloClusterName = d.ApolloClusterName
	r.ApolloNamespace = d.ApolloNamespace
	r.JsVersion = d.JsVersion
	r.ModelBranch = d.ModelBranch
	r.DeployEnv = env.Name
	r.DeployEnvStatus = env.Status
	r.Replics = project.CopyCount
	r.ContainerPort = project.ContainerPort
	r.ServiceListenPort = project.ServiceListenPort
	r.CpuRequest = strconv.Itoa(project.CpuMinRequire) + "m"
	r.CpuLimit = strconv.Itoa(project.CpuMaxRequire) + "m"
	r.MemoryRequest = strconv.Itoa(project.MemoryMinRequire) + "Mi"
	r.MemoryLimit = strconv.Itoa(project.MemoryMaxRequire) + "Mi"
	r.IfStorageLocale = *project.IfStorageLocale
	r.StoragePath = project.StoragePath
	r.IfCheckPodsStatus = *project.IfCheckPodsStatus
	r.IfUseIstio = *project.IfUseIstio
	r.IfUseApolloOfflineEnv = *project.IfUseApolloOfflineEnv

	if r.ProductName == "default" {
		r.ProductName = project.OwnedProduct
	}

	return r
}

func (r *ReqJenkinsBuild) SetReplics(envName string, podsNumString string) *ReqJenkinsBuild {
	podsNumList := strings.Split(podsNumString, ",")
	log.Info("ttttttttt : " + podsNumString)
	if len(podsNumList) <= 2 {
		r.Replics = 1
		return r
	}

	for _, val := range podsNumList {
		log.Info(val)
		tmpList := strings.Split(val, ":")
		log.Info(tmpList, envName)
		if len(tmpList) > 1 && envName == tmpList[0] {
			r.Replics = tools.StringToInt(tmpList[1])
		}
	}
	return r
}

func (r *ReqJenkinsBuild) SetUnityAppName(unityAppId int) {
	r.UnityAppName = "no_unity"
	if r.CodeLanguage == "android" && r.IfAddUnityProject {
		var unityAppName struct {
			Name string `json:"name"`
		}
		db := DB
		db = db.Table(Project{}.TableName())
		db = db.Select("name")
		db = db.Where("id = ?", unityAppId).First(&unityAppName)
		log.Info(unityAppName)
		r.UnityAppName = unityAppName.Name
	}
	log.Info(fmt.Sprintf("unity app name : %s", r.UnityAppName))
}

type JenkinsBuildResponse struct {
	Info   string `json:"info"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

type PostChange struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}

type ServiceCallJenkinsTriggerRespone struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Res  string `json:"res"`
	Data struct {
		Id        int64  `json:"id"`
		Triggered bool   `json:"triggered"`
		Url       string `json:"url"`
	} `json:"data"`
}
type ServiceCallJenkinsJobUpdateRespone struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Res  string `json:"res"`
}
