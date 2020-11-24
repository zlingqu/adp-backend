package model

import "time"

type ProjectIdName struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"type:varchar(80);unique_index"`
}

type ProjectIdNameGit struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Name          string `json:"name" gorm:"type:varchar(80);unique_index"`
	GitRepository string `json:"git_repository" gorm:"type:varchar(300)"`
}

type ProjectIdNameGitLang struct {
	ID                   uint   `json:"id" gorm:"primary_key"`
	Name                 string `json:"name" gorm:"type:varchar(80);unique_index"`
	GitRepository        string `json:"git_repository" gorm:"type:varchar(300)"`
	LanguageType         string `json:"language_type" gorm:"type:varchar(80)"`
	IfUseModel           bool   `json:"if_use_model"`
	IfUseGitManagerModel bool   `json:"if_use_git_manager_model"`
	ModelGitRepository   string `json:"model_git_repository"  gorm:"type:varchar(300)"`
}

type ProjectIdNameGitLangProduct struct {
	ID                   uint   `json:"id" gorm:"primary_key"`
	Name                 string `json:"name" gorm:"type:varchar(80);unique_index"`
	GitRepository        string `json:"git_repository" gorm:"type:varchar(300)"`
	LanguageType         string `json:"language_type" gorm:"type:varchar(80)"`
	OwnedProduct         string `json:"owned_product" gorm:"type:varchar(80)"`
	IfUseModel           bool   `json:"if_use_model"`
	IfUseGitManagerModel bool   `json:"if_use_git_manager_model"`
	ModelGitRepository   string `json:"model_git_repository"  gorm:"type:varchar(300)"`
}

type Project struct {
	ID                        uint      `json:"id" gorm:"primary_key"`
	CreatedAt                 time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	UpdatedAt                 time.Time `json:"updated_at"`
	Name                      string    `json:"name" gorm:"type:varchar(80);unique_index"`
	Status                    string    `json:"status" gorm:"type:varchar(80)"`
	OwnedProduct              string    `json:"owned_product" gorm:"type:varchar(80)"` ////
	Description               string    `json:"description" gorm:"type:text"`          ////
	GitRepository             string    `json:"git_repository" gorm:"type:varchar(300)"`
	LanguageType              string    `json:"language_type" gorm:"type:varchar(80)"`
	IfAddUnityProject         bool     `json:"if_add_unity_project"`
	UnityAppId                int       `json:"unity_app_id"`
	DeployEnvType             string    `json:"deploy_env_type" gorm:"type:varchar(80)"`
	ServeType                 string    `json:"serve_type" gorm:"type:varchar(80)"`
	ReplicationControllerType string    `json:"replication_controller_type" gorm:"type:varchar(80)"`
	IfUseDomainName           bool      `json:"if_use_domain_name"`
	DomainName                string    `json:"domain_name" gorm:"type:varchar(80)"`
	IfUseHttps                bool     `json:"if_use_https"`
	IfUseHttp                 bool     `json:"if_use_http"`
	IfUseGrpc                 bool     `json:"if_use_grpc"`
	IfUseGbs                  bool     `json:"if_use_gbs"`
	IfUseSticky               bool     `json:"if_use_sticky"`
	IfCompile                 bool     `json:"if_compile"`
	IfCompileCache            bool     `json:"if_compile_cache"`
	IfCompileParam            bool     `json:"if_compile_param"`
	CompileParam              string    `json:"compile_param" gorm:"type:varchar(300)"`
	IfCompileImage            bool     `json:"if_compile_image"`
	CompileImage              string    `json:"compile_image" gorm:"type:varchar(500)"`
	IfMakeImage               bool     `json:"if_make_image"`
	IfUseModel                bool     `json:"if_use_model"`
	IfUseGitManagerModel      bool     `json:"if_use_git_manager_model"`
	ModelGitRepository        string    `json:"model_git_repository"  gorm:"type:varchar(300)"`
	IfSaveModelBuildComputer  bool     `json:"if_save_model_build_computer"`
	IfUseConfigmap            bool     `json:"if_use_configmap"`
	IfUseAutoDeployFile       bool     `json:"if_use_auto_deploy_file"`
	AutoDeployContent         string    `json:"auto_deploy_content" gorm:"type:text"`
	IfDeploy                  bool     `json:"if_deploy"`
	PodsNum                   string    `json:"pods_num"`
	CopyCount                 int       `json:"copy_count"`
	ContainerPort             int       `json:"container_port"`
	ServiceListenPort         string    `json:"service_listen_port" gorm:"type:varchar(80)"`
	IfUseCustomDockerfile     bool     `json:"if_use_custom_dockerfile"`
	IfUseRootDockerfile       bool     `json:"if_use_root_dockerfile"`
	DockerfileContent         string    `json:"dockerfile_content" gorm:"type:text"`
	CpuMinRequire             int       `json:"cpu_min_require"`
	CpuMaxRequire             int       `json:"cpu_max_require"`
	MemoryMinRequire          int       `json:"memory_min_require"`
	MemoryMaxRequire          int       `json:"memory_max_require"`
	IfStorageLocale           bool     `json:"if_storage_locale"`
	StoragePath               string    `json:"storage_path" gorm:"type:varchar(80)"`
	IfUseGpuCard              bool     `json:"if_use_gpu_card"`
	GpuControlMode            string    `json:"gpu_control_mode" gorm:"type:varchar(80)"`
	GpuCardCount              int       `json:"gpu_card_count"`
	GpuMemCount               int       `json:"gpu_mem_count"`
	IfUseIstio                bool     `json:"if_use_istio"`
	IfUseApolloOfflineEnv     bool     `json:"if_use_apollo_offline_env"`
	IfNeedCheck               bool     `json:"if_need_check"`
	IfCheckPodsStatus         bool     `json:"if_check_pods_status"`
}

func (Project) TableName() string {
	return "project"
}

func NewProject() *Project {
	p := &Project{}
	// if !Model.HasTable(p.TableName()) {
	// 	Model.CreateTable(p)
	// }
	if Model.HasTable(p.TableName()) { //判断表是否存在
		Model.AutoMigrate(p) //存在就自动适配表，也就说原先没字段的就增加字段
	} else {
		Model.CreateTable(p) //不存在就创建新表
	}
	return p
}

type GetProject struct {
	Name string `form:"name"`
	Page int64  `form:"page"`
	Size int64  `form:"size"`
}

type GetByID struct {
	ID string `uri:"id" binding:"required"`
}

type PostIds struct {
	Ids []string `json:"ids"`
}

type JenkinsJob struct {
	GitAddress  string `json:"git_address"`
	AppName     string `json:"app_name"`
	ProductName string `json:"product_name"`
	Action      string `json:"action"`
}

type JenkinsResponse struct {
	Info   string `json:"info"`
	Status string `json:"status"`
}
