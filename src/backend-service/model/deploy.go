package model

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/zuoshenglo/libs/logs/logrus"
)

type UpdateDeploy struct {
	AppId                        uint   `json:"app_id" gorm:"comment:应用的ID"`
	Branch                       string `json:"branch" gorm:"type:varchar(80);comment:git仓库的分支名"`
	EnvId                        uint   `json:"env_id" gorm:"comment:环境ID"`
	GitCommitId                  string `json:"git_commit_id" gorm:"type:varchar(150);default:last;comment:git仓库的commitID"`
	GitTag                       string `json:"git_tag" gorm:"type:varchar(150);comment:git仓库的tag"`
	ID                           uint   `json:"id" gorm:"primary_key;comment:工单ID"`
	Name                         string `json:"name" gorm:"type:varchar(80);comment:工单名字"`
	OwnerEnglishName             string `json:"owner_english_name" gorm:"type:varchar(80);comment:工单所属英文名"`
	OwnerChinaName               string `json:"owner_china_name" gorm:"type:varchar(80);comment:工单所属中文名"`
	Status                       string `json:"status" gorm:"type:varchar(20);comment:工单状态"`
	VersionControlMode           string `json:"version_control_mode" gorm:"type:varchar(80);comment:代码版本控制方式"`
	PodNums                      int    `json:"pod_nums" gorm:"type:int(11);default:1;comment:POD数量"`
	IfStorageLocale              *bool  `json:"if_storage_locale" gorm:"default:0;comment:是否需要存储"`
	StoragePath                  string `json:"storage_path" gorm:"type:varchar(512);comment:存储路径"`
	CpuMinRequire                *int   `json:"cpu_min_require" gorm:"type:int(11);default:'100';comment:CPU需求最小值"`
	CpuMaxRequire                *int   `json:"cpu_max_require" gorm:"type:int(11);default:200;comment:CPU最大限制"`
	MemoryMinRequire             *int   `json:"memory_min_require" gorm:"type:int(11);default:200;comment:内存需求最小值"`
	MemoryMaxRequire             *int   `json:"memory_max_require" gorm:"type:int(11);default:400;comment:内存最大限制"`
	GpuControlMode               string `json:"gpu_control_mode" gorm:"type:varchar(80);default:'mem';comment:gpu使用方式"`
	GpuCardCount                 int    `json:"gpu_card_count" gorm:"type:int(11);default:1;comment:gpu卡数量"`
	GpuMemCount                  int    `json:"gpu_mem_count" gorm:"type:int(11);default:2;comment:gpu显存大小"`
	GpuType                      string `json:"gpu_type" gorm:"type:varchar(512);default:'all';comment:gpu型号"`
	IfUsePodAntiAffinity         *bool  `json:"if_use_pod_anti_affinity" gorm:"default:0;comment:是否开启pod反亲和"`
	IfUseApollo                  *bool  `json:"if_use_apollo" gorm:"default:1;comment:是否需要使用apollo配置中心"`
	ApolloClusterName            string `json:"apollo_cluster_name" gorm:"type:varchar(80);comment:apollo集群名字"`
	ApolloNamespace              string `json:"apollo_namespace" gorm:"type:varchar(80);comment:apollo的namespace"`
	IfUseApolloForDockerfile     *bool  `json:"if_use_apollo_for_dockerfile" gorm:"default:0;comment:是否将apollo配置注入环境变量到Dockerfile"`
	ApolloEnvForDockerfile       string `json:"apollo_env_for_dockerfile" gorm:"default:'prd';type:varchar(80);comment:关联的apollo环境名字"`
	ApolloClusterForDockerfile   string `json:"apollo_cluster_for_dockerfile" gorm:"default:'default';type:varchar(80);comment:关联的apollo集群名字"`
	ApolloNamespaceForDockerfile string `json:"apollo_namespace_for_dockerfile" gorm:"default:'application';type:varchar(80);comment:关联的apollo空间名"`
	K8sNamespace                 string `json:"k8s_namespace" gorm:"type:varchar(80);comment:k8s的namespace"`
	NodeSpecialFor               string `json:"node_special_for" gorm:"type:varchar(20);default:'None';comment:特殊node标识"`
	YamlEnv                      string `json:"yaml_env" gorm:"type:varchar(1024);default:None;comment:yaml文件需要注入的环境变量"`
	AndroidFlavor                string `json:"android_flavor" gorm:"type:varchar(80);default:default;comment:安卓编译渠道号"`
	DomainBefore                 string `json:"domain_before" gorm:"type:varchar(20);default:None;comment:域名前缀值"`
	DomainMiddle                 string `json:"domain_middle" gorm:"type:varchar(80);default:None;comment:域名中间值"`
	DomainAfter                  string `json:"domain_after" gorm:"type:varchar(80);default:None;comment:域名后缀值"`
	DomainPath                   string `json:"domain_path" gorm:"type:varchar(256);default:None;comment:域名path路径"`
	JsVersion                    string `json:"js_version" gorm:"type:varchar(80)"`
}

type Deploy struct {
	UpdateDeploy
	CreatedAt         time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastDeploy        time.Time `json:"last_deploy"`
	LastBuildInfo     string    `json:"last_build_info" gorm:"type:text"`
	JenkinsBuildToken string    `json:"jenkins_build_token" gorm:"type:text"`
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
	GitAddress                   string `json:"git_address"`
	AppName                      string `json:"app_name"`
	ProductName                  string `json:"product_name"`
	CodeLanguage                 string `json:"code_language"`
	IfAddUnityProject            bool   `json:"if_add_unity_project"`
	UnityAppName                 string `json:"unity_app_name"`
	DeployEnvType                string `json:"deploy_env_type"`
	IfCompile                    bool   `json:"if_compile"`
	IfMakeImage                  bool   `json:"if_make_image"`
	IfUseGrpc                    bool   `json:"if_use_grpc"`
	IfUseGbs                     bool   `json:"if_use_gbs"`
	IfUseSticky                  bool   `json:"if_use_sticky"`
	IfDeploy                     bool   `json:"if_deploy"`
	IfUseModel                   bool   `json:"if_use_model"`
	IfUseGitManagerModel         bool   `json:"if_use_git_manager_model"`
	ModelGitRepository           string `json:"model_git_repository"`
	IfSaveModelBuildComputer     bool   `json:"if_save_model_build_computer"`
	IfUseAutoDeployFile          bool   `json:"if_use_auto_deploy_file"`
	AutoDeployContent            string `json:"auto_deploy_content"`
	IfUseCustomDockerfile        bool   `json:"if_use_custom_dockerfile"`
	IfUseRootDockerfile          bool   `json:"if_use_root_dockerfile"`
	DockerfileContent            string `json:"dockerfile_content"`
	ServeType                    string `json:"serve_type"`
	ReplicationControllerType    string `json:"replication_controller_type"`
	IfUseGpuCard                 bool   `json:"if_use_gpu_card"`
	GpuControlMode               string `json:"gpu_control_mode"`
	GpuCardCount                 int    `json:"gpu_card_count"`
	GpuMemCount                  int    `json:"gpu_mem_count"`
	GpuType                      string `json:"gpu_type"`
	BranchName                   string `json:"branch_name"`
	Version                      string `json:"version"`
	VersionControlMode           string `json:"version_control_mode"`
	GitCommitId                  string `json:"git_commit_id"`
	GitTag                       string `json:"git_tag"`
	IfUsePodAntiAffinity         bool   `json:"if_use_pod_anti_affinity"`
	IfUseApollo                  bool   `json:"if_use_apollo"`
	ApolloClusterName            string `json:"apollo_cluster_name"`
	ApolloNamespace              string `json:"apollo_namespace"`
	DeployEnv                    string `json:"deploy_env"`
	DeployEnvStatus              string `json:"deploy_env_status"`
	Replics                      int    `json:"replics"`
	ContainerPort                int    `json:"container_port"`
	ServiceListenPort            string `json:"service_listen_port"`
	CpuRequest                   string `json:"cpu_request"`
	CpuLimit                     string `json:"cpu_limit"`
	MemoryRequest                string `json:"memory_request"`
	MemoryLimit                  string `json:"memory_limit"`
	IfStorageLocale              bool   `json:"if_storage_locale"`
	StoragePath                  string `json:"storage_path"`
	IfCheckPodsStatus            bool   `json:"if_check_pods_status"`
	IfUseIstio                   bool   `json:"if_use_istio"`
	JsVersion                    string `json:"js_version"`
	YamlEnv                      string `json:"yaml_env"`
	NodeSpecialFor               string `json:"node_special_for"`
	AndroidFlavor                string `json:"android_flavor"`
	IfUseApolloForDockerfile     bool   `json:"if_use_apollo_for_dockerfile"`
	ApolloEnvForDockerfile       string `json:"apollo_env_for_dockerfile"`
	ApolloClusterForDockerfile   string `json:"apollo_cluster_for_dockerfile"`
	ApolloNamespaceForDockerfile string `json:"apollo_namespace_for_dockerfile"`
	DomainBefore                 string `json:"domain_before"`
	DomainMiddle                 string `json:"domain_middle"`
	DomainAfter                  string `json:"domain_after"`
	DomainPath                   string `json:"domain_path"`
}

func (r *ReqJenkinsBuild) SetReqJenkinsBuildData(env Env, project Project, d Deploy) *ReqJenkinsBuild {
	r.GitAddress = project.GitRepository
	r.AppName = project.Name
	r.ProductName = d.K8sNamespace
	r.CodeLanguage = project.LanguageType
	r.IfAddUnityProject = *project.IfAddUnityProject
	r.DeployEnvType = project.DeployEnvType
	r.IfCompile = *project.IfCompile
	r.IfMakeImage = *project.IfMakeImage
	r.IfUseGrpc = *project.IfUseGrpc
	r.IfUseGbs = *project.IfUseGbs
	r.IfUseSticky = *project.IfUseSticky
	r.IfDeploy = *project.IfDeploy
	r.IfUseModel = *project.IfUseModel
	r.IfUseGitManagerModel = *project.IfUseGitManagerModel
	r.ModelGitRepository = project.ModelGitRepository
	r.IfSaveModelBuildComputer = *project.IfSaveModelBuildComputer
	r.IfUseAutoDeployFile = *project.IfUseAutoDeployFile
	r.AutoDeployContent = project.AutoDeployContent
	r.IfUseCustomDockerfile = *project.IfUseCustomDockerfile
	r.IfUseRootDockerfile = *project.IfUseRootDockerfile
	r.DockerfileContent = project.DockerfileContent
	r.ServeType = project.ServeType
	r.ReplicationControllerType = project.ReplicationControllerType
	r.GpuControlMode = d.GpuControlMode
	r.GpuCardCount = d.GpuCardCount
	r.GpuMemCount = d.GpuMemCount
	r.GpuType = d.GpuType
	r.BranchName = d.Branch
	r.VersionControlMode = d.VersionControlMode
	r.GitCommitId = d.GitCommitId
	r.GitTag = d.GitTag
	r.IfUsePodAntiAffinity = *d.IfUsePodAntiAffinity
	r.IfUseApollo = *d.IfUseApollo
	r.ApolloClusterName = d.ApolloClusterName
	r.ApolloNamespace = d.ApolloNamespace
	r.JsVersion = d.JsVersion
	r.YamlEnv = d.YamlEnv
	r.NodeSpecialFor = d.NodeSpecialFor
	r.DeployEnv = env.Name
	r.DeployEnvStatus = env.Status
	r.Replics = d.PodNums
	r.ContainerPort = project.ContainerPort
	r.ServiceListenPort = project.ServiceListenPort
	r.CpuRequest = strconv.Itoa(*d.CpuMinRequire) + "m"
	r.CpuLimit = strconv.Itoa(*d.CpuMaxRequire) + "m"
	r.MemoryRequest = strconv.Itoa(*d.MemoryMinRequire) + "Mi"
	r.MemoryLimit = strconv.Itoa(*d.MemoryMaxRequire) + "Mi"
	r.IfStorageLocale = *d.IfStorageLocale
	r.StoragePath = d.StoragePath
	r.IfCheckPodsStatus = *project.IfCheckPodsStatus
	r.IfUseIstio = *project.IfUseIstio
	r.AndroidFlavor = d.AndroidFlavor
	r.IfUseApolloForDockerfile = *d.IfUseApolloForDockerfile
	r.ApolloEnvForDockerfile = d.ApolloEnvForDockerfile
	r.ApolloClusterForDockerfile = d.ApolloClusterForDockerfile
	r.ApolloNamespaceForDockerfile = d.ApolloNamespaceForDockerfile
	r.DomainBefore = d.DomainBefore
	r.DomainMiddle = d.DomainMiddle
	r.DomainAfter = d.DomainAfter
	r.DomainPath = d.DomainPath

	if r.ProductName == "default" {
		r.ProductName = project.OwnedProduct
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
