package handler

import (
	"app-deploy-platform/3rd-api/jenkins/config"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	glibs "github.com/zuoshenglo/go-base-libs"
	"github.com/zuoshenglo/tools"
)

type JenkinsBuild struct {
	DmaiBaseDevopsHttp
	JenkinsJobBuild
	JenkinsLastBuild
	parameter map[string][]map[string]interface{}
}

//return &JenkinsBuild{
//	parameter:          make(map[string][]map[string]interface{}, 0),
//	DmaiBaseDevopsHttp: DmaiBaseDevopsHttp{},
//	JenkinsJobBuild:    JenkinsJobBuild{},
//	JenkinsLastBuild:   JenkinsLastBuild{},
//}

func NewJenkinsBuild(c *gin.Context) {
	returnDate := gin.H{
		"status": "faild",
		"info":   "init",
		"url":    "kongkong",
	}

	j := JenkinsBuild{
		parameter:          make(map[string][]map[string]interface{}, 0),
		DmaiBaseDevopsHttp: DmaiBaseDevopsHttp{},
		JenkinsJobBuild:    JenkinsJobBuild{},
		JenkinsLastBuild:   JenkinsLastBuild{},
	}

	// user req param

	var userJson map[string]interface{}

	if err := c.ShouldBindJSON(&userJson); err != nil {
		log.Error(err)
		returnDate["info"] = "Request parameter error"
		c.JSON(http.StatusOK, returnDate)
		return
	}

	log.Info(userJson["app_name"])
	log.Info(userJson["branch_name"])

	// init j.addParameter
	j.parameter = make(map[string][]map[string]interface{}, 0)
	// jenkins build param
	//j.addParameter("NODE_ENV", strings.Replace(userJson["deploy_env"].(string), "prd", "prod", -1))
	j.addParameter("NODE_ENV", "dev")
	j.addParameter("DEPLOY_ENV", "dev")
	j.addParameter("APP_NAME", userJson["app_name"].(string))
	j.addParameter("ENV_TYPE", userJson["deploy_env_type"].(string))
	j.addParameter("REPLICAS", userJson["replics"].(float64))
	//j.addParameter("CONTAINER_PORT", userJson["container_port"].(float64))
	j.addParameter("CONTAINER_PORT", userJson["service_listen_port"].(string))
	j.addParameter("CPU_REQUEST", userJson["cpu_request"].(string))
	j.addParameter("CPU_LIMIT", userJson["cpu_limit"].(string))
	j.addParameter("MEMORY_REQUEST", userJson["memory_request"].(string))
	j.addParameter("MEMORY_LIMIT", userJson["memory_limit"].(string))
	j.addParameter("NAMESPACE", userJson["product_name"].(string))
	j.addParameter("GIT_ADDRESS", userJson["git_address"].(string))
	j.addParameter("CODE_LANGUAGE", userJson["code_language"].(string))
	j.addParameter("COMPILE", userJson["if_compile"].(bool))
	j.addParameter("COMPILE_PARAM", userJson["compile_param"].(string))
	j.addParameter("IF_MAKE_IMAGE", userJson["if_make_image"].(bool))
	j.addParameter("DEPLOY", userJson["if_deploy"].(bool))
	j.addParameter("DOMAIN", userJson["domain_name"].(string))
	j.addParameter("IF_USE_HTTPS", userJson["if_use_https"].(bool))
	j.addParameter("IF_USE_HTTP", userJson["if_use_http"].(bool))
	j.addParameter("CUSTOM_KUBERNETES_DEPLOY_TEMPLATE", userJson["if_use_auto_deploy_file"].(bool))
	j.addParameter("CUSTOM_KUBERNETES_DEPLOY_TEMPLATE_CONTENT", userJson["auto_deploy_content"].(string))
	j.addParameter("CUSTOM_DOCKERFILE", userJson["if_use_custom_dockerfile"].(bool))
	j.addParameter("CUSTOM_DOCKERFILE_CONTENT", userJson["dockerfile_content"].(string))
	j.addParameter("SERVICE_TYPE", userJson["serve_type"].(string))
	//j.addParameter("ENV_TYPE", j.getGpuType(userJson["if_use_gpu_card"].(bool)))
	j.addParameter("GPU_CONTROL_MODE", userJson["gpu_control_mode"].(string))
	j.addParameter("GPU_CARD_COUNT", userJson["gpu_card_count"].(float64))
	j.addParameter("GPU_MEM_COUNT", userJson["gpu_mem_count"].(float64))
	j.addParameter("GPU_TYPE", userJson["gpu_type"].(string))
	j.addParameter("GIT_VERSION", userJson["git_commit_id"])
	j.addParameter("GIT_TAG", userJson["git_tag"])
	j.addParameter("APOLLO_CLUSTER_NAME", userJson["apollo_cluster_name"])
	j.addParameter("APOLLO_NAMESPACE", userJson["apollo_namespace"])
	j.addParameter("JS_VERSION", userJson["js_version"])
	j.addParameter("BRANCH_NAME", userJson["branch_name"])
	j.addParameter("VERSION_CONTROL_MODE", userJson["version_control_mode"])
	j.addParameter("USE_MODEL", userJson["if_use_model"].(bool))
	j.addParameter("USE_CONFIGMAP", userJson["if_use_configmap"].(bool))
	j.addParameter("IF_CHECK_PODS_STATUS", userJson["if_check_pods_status"].(bool))
	j.addParameter("IF_STORAGE_LOCALE", userJson["if_storage_locale"].(bool))
	j.addParameter("STORAGE_PATH", userJson["storage_path"].(string))
	j.addParameter("DEPLOY_MASTER_PASSWORD", "dmai2019999")
	j.addParameter("USE_SERVICE", true)
	j.addParameter("BUILD_PLATFORM", "adp")
	// all
	j.addParameter("GLOABL_STRING", tools.BoolToString(userJson["if_use_grpc"].(bool))+":::"+
		tools.BoolToString(userJson["if_use_sticky"].(bool))+":::"+
		userJson["replication_controller_type"].(string)+":::"+
		tools.BoolToString(userJson["if_add_unity_project"].(bool))+":::"+
		userJson["unity_app_name"].(string)+":::"+
		tools.BoolToString(userJson["if_use_root_dockerfile"].(bool))+":::"+
		tools.BoolToString(userJson["if_compile_param"].(bool))+":::"+
		tools.BoolToString(userJson["if_compile_image"].(bool))+":::"+
		userJson["compile_image"].(string)+":::"+
		tools.BoolToString(userJson["if_use_gbs"].(bool))+":::"+
		tools.BoolToString(userJson["if_compile_cache"].(bool))+":::"+
		tools.BoolToString(userJson["if_use_model"].(bool))+":::"+
		tools.BoolToString(userJson["if_use_git_manager_model"].(bool))+":::"+
		tools.BoolToString(userJson["if_save_model_build_computer"].(bool))+":::"+
		userJson["model_git_repository"].(string)+":::"+
		userJson["model_branch"].(string)+":::"+
		userJson["deploy_env_status"].(string)+":::"+
		userJson["deploy_env"].(string)+":::"+
		tools.BoolToString(userJson["if_use_istio"].(bool))+":::"+
		tools.BoolToString(userJson["if_use_apollo_offline_env"].(bool)))

	//
	data, err := json.Marshal(j.parameter)
	if err != nil {
		log.Error(err)
		returnDate["info"] = "json marshal faild"
		c.JSON(http.StatusOK, returnDate)
		return
	}

	log.Info(string(data))
	//request url
	jenkinsBaseUrl := config.GetEnv().JenkinsAddress + "/job/" + userJson["app_name"].(string) + "/job/" + url.QueryEscape(userJson["branch_name"].(string))
	log.Info(jenkinsBaseUrl)

	// request jenkins
	req := glibs.NewHttpRequestCustom([]byte(""), "POST", jenkinsBaseUrl+"/build").SetRequestProtocol("http").SetContentType("application/x-www-form-urlencoded").SetFormKeyValues("json", string(data))
	req.SetBasicAuth(config.GetEnv().JenkinsUser, config.GetEnv().JenkinsPasswd)
	result, err := req.ExecRequest()
	if err != nil {
		log.Error(err)
		returnDate["info"] = "第一次请求jenkins构建，执行请求失败！"
		c.JSON(http.StatusOK, returnDate)
		return
	}

	log.Info(result)

	// check return
	if result != "" {
		log.Warn("第一次请求Jenkins进行构建失败，可能是jenkins提供的参数不匹配，进行参数化更新")

		// 参数化更新
		//j.addParameter("GIT_VERSION", "update")
		//data, _ := json.Marshal(j.parameter)
		//log.Info(string(data))
		updateString := "{\"parameter\":[{\"name\":\"GIT_VERSION\",\"value\":\"update\"}]}"
		req := glibs.NewHttpRequestCustom([]byte(""), "POST", jenkinsBaseUrl+"/build").SetRequestProtocol("http").SetContentType("application/x-www-form-urlencoded").SetFormKeyValues("json", updateString)
		req.SetBasicAuth(config.GetEnv().JenkinsUser, config.GetEnv().JenkinsPasswd)
		result, err := req.ExecRequest()
		if err != nil {
			log.Error(err)
			returnDate["info"] = "第二次请求jenkins构建，参数化配置更新，执行请求失败！"
			c.JSON(http.StatusOK, returnDate)
			return
		}

		if result != "" {
			log.Error("第二次请求jenkins构建，参数化配置更新，jenkins失败！")
			returnDate["info"] = "第二次请求jenkins构建，参数化配置更新，jenkins失败！"
			c.JSON(http.StatusOK, returnDate)
			return
		}

		log.Info("第二次请求jenkins构建，参数化配置更新，jenkins成功，再次请求构建！")

		time.Sleep(30 * time.Second)

		req = glibs.NewHttpRequestCustom([]byte(""), "POST", jenkinsBaseUrl+"/build").SetRequestProtocol("http").SetContentType("application/x-www-form-urlencoded").SetFormKeyValues("json", string(data))
		req.SetBasicAuth(config.GetEnv().JenkinsUser, config.GetEnv().JenkinsPasswd)
		result, err = req.ExecRequest()

		if err != nil {
			log.Error("第三次请求Jenkins进行构建， 执行请求失败")
			returnDate["info"] = "第三次请求Jenkins进行构建， 执行请求失败"
			c.JSON(http.StatusOK, returnDate)
			return
		}

		if result != "" {
			log.Error("第三次请求jenkins构建，jenkins失败！")
			returnDate["info"] = "第三次请求jenkins构建，jenkins失败！"
			c.JSON(http.StatusOK, returnDate)
			return
		}
		log.Info("请求jenkins进行构建，执行成功！")
	}

	// get lastBuild
	log.Info("begin get last build")
	jenkinsdLastBuildApiUrl := jenkinsBaseUrl + "/lastBuild/api/json"
	log.Info(jenkinsdLastBuildApiUrl)
	log.Info(config.GetEnv().JenkinsUser)
	reqLast := glibs.NewHttpRequestCustom([]byte(""), "POST", jenkinsdLastBuildApiUrl).SetRequestProtocol("http").SetContentType("")
	reqLast.SetBasicAuth(config.GetEnv().JenkinsUser, config.GetEnv().JenkinsPasswd)
	resultLast, errLast := reqLast.ExecRequest()
	if errLast != nil {
		log.Error(errLast)
	}

	// Unmarshal
	log.Info(resultLast)
	jsonLastErr := json.Unmarshal([]byte(resultLast), &j.JenkinsLastBuild)
	if jsonLastErr != nil {
		log.Error(jsonLastErr)
		returnDate["info"] = fmt.Sprintf("%s", jsonLastErr)
		c.JSON(http.StatusOK, returnDate)
		return
	}

	log.Info(j.JenkinsLastBuild)
	buildIdInt, _ := strconv.Atoi(j.JenkinsLastBuild.Id)
	buildId := fmt.Sprintf("%d", buildIdInt+1)
	returnDate["status"] = "ok"
	returnDate["info"] = "build success"
	returnDate["url"] = config.GetEnv().JenkinsPipelineURL + userJson["app_name"].(string) + "/detail/" + url.QueryEscape(userJson["branch_name"].(string)) + "/" + buildId + "/pipeline"

	log.Info(returnDate)
	c.JSON(http.StatusOK, returnDate)
	return
}

func (j *JenkinsBuild) addParameter(name string, value interface{}) {
	j.parameter["parameter"] = append(j.parameter["parameter"], map[string]interface{}{"name": name, "value": value})
	log.Info(j.parameter)
	return
}

func (j *JenkinsBuild) getGpuType(useGpu bool) string {
	if useGpu {
		return "gpu"
	}
	return "cpu"
}

type JenkinsJob struct {
	GitAddress    string `json:"git_address"`
	JobName       string `json:"app_name"`
	ProductName   string `json:"product_name"`
	ConfigXmlPath string `json:"config_xml_path"`
}

type JenkinsJobBuild struct {
	JobName    string `json:"app_name"`
	BranchName string `json:"branch_name"`
}

type JenkinsLastBuild struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func JenkinsJobCfgFile(appName string, gitAddress string) string {
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

/*
1. 对请求的方法进行封装。
2. 对统一的json数据进行抽取。
3. 统一的取值方式，针对多层嵌套的json格式的数据。
*/

type DmaiBaseDevopsHttp struct {
	rep        http.ResponseWriter
	req        *http.Request
	UserJson   string
	ResponJson string
}
