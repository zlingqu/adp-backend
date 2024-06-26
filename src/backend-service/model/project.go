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
	IfDeploy             bool   `json:"if_deploy"`
	IfUseAutoDeployFile  bool   `json:"if_use_auto_deploy_file"`
	DeployEnvType        string `json:"deploy_env_type"`
	ServeType            string `json:"serve_type"`
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
	IfAddUnityProject         *bool     `json:"if_add_unity_project" gorm:"default:0"`
	UnityAppId                int       `json:"unity_app_id"`
	DeployEnvType             string    `json:"deploy_env_type" gorm:"type:varchar(80)"`
	ServeType                 string    `json:"serve_type" gorm:"type:varchar(80)"`
	ReplicationControllerType string    `json:"replication_controller_type" gorm:"type:varchar(80)"`
	IfUseGrpc                 *bool     `json:"if_use_grpc" gorm:"default:0"`
	IfUseGbs                  *bool     `json:"if_use_gbs" gorm:"default:0"`
	IfUseSticky               *bool     `json:"if_use_sticky" gorm:"default:0"`
	IfCompile                 *bool     `json:"if_compile" gorm:"default:0"`
	IfMakeImage               *bool     `json:"if_make_image" gorm:"default:0"`
	IfUseModel                *bool     `json:"if_use_model" gorm:"default:0"`
	IfUseGitManagerModel      *bool     `json:"if_use_git_manager_model" gorm:"default:0"`
	ModelGitRepository        string    `json:"model_git_repository"  gorm:"type:varchar(300)"`
	IfSaveModelBuildComputer  *bool     `json:"if_save_model_build_computer" gorm:"default:0"`
	IfUseAutoDeployFile       *bool     `json:"if_use_auto_deploy_file" gorm:"default:0"`
	IfDeploy                  *bool     `json:"if_deploy" gorm:"default:0"`
	ContainerPort             int       `json:"container_port"`
	ServiceListenPort         string    `json:"service_listen_port" gorm:"type:varchar(80)"`
	IfUseCustomDockerfile     *bool     `json:"if_use_custom_dockerfile" gorm:"default:0"`
	IfUseRootDockerfile       *bool     `json:"if_use_root_dockerfile" gorm:"default:0"`
	DockerfileContent         string    `json:"dockerfile_content" gorm:"type:text"`
	IfUseIstio                *bool     `json:"if_use_istio" gorm:"default:0"`
	IfNeedCheck               *bool     `json:"if_need_check" gorm:"default:0"`
	IfCheckPodsStatus         *bool     `json:"if_check_pods_status" gorm:"default:0"`
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
