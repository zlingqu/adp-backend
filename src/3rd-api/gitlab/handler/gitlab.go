package handler

import (
	"app-deploy-platform/3rd-api/gitlab/config"
	svc "app-deploy-platform/3rd-api/gitlab/service"
	"app-deploy-platform/common/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// func ProjectBranchs(c *gin.Context) {

// 	gitHttpAddress := c.DefaultQuery("http_url_to_repo", "error")

// 	if gitHttpAddress == "error" {
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": 0,
// 			"res":  "error",
// 			"msg":  "http_url_to_repo参数错误。",
// 		})
// 		return
// 	}

// 	// 查询项目id
// 	gitLab := NewGitLab().GetProjectName(gitHttpAddress).GetProjectId().GetProjectBranchs()

// 	c.JSON(http.StatusOK, gin.H{
// 		"code": 0,
// 		"res":  gitLab.Res,
// 		"msg":  gitLab.Msg,
// 		"data": gitLab.ProjectBranches,
// 	})
// }

func GetCommitID(c *gin.Context) {

	gitHttpUrl := c.DefaultQuery("http_url_to_repo", "error")
	branchName := c.DefaultQuery("branch", "error")
	commitID := svc.GetCommitIDByRepourlAndBranch(gitHttpUrl, branchName)
	fmt.Println(gitHttpUrl, branchName)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "OK",
		"msg":  "获取成功",
		"data": commitID,
		// "data": 2000,
	})
}

type GitLab struct {
	Id                     int64                    `json:"id"`
	ProjectName            string                   `json:"project_name"`
	ProjectGitAddress      string                   `json:"project_git_address"`
	ProjectBranches        []string                 `json:"project_branchs"`
	GitPrivateToken        string                   `json:"git_private_token"`
	GitProjectApiAddress   string                   `json:"git_project_api_address"`
	Res                    string                   `json:"res"`
	Msg                    string                   `json:"msg"`
	GitGetIdResponse       []map[string]interface{} `json:"git_get_id_response"`
	GitGetBranchesResponse []map[string]interface{} `json:"git_get_branches_response"`
}

func NewGitLab() *GitLab {
	return &GitLab{
		GitPrivateToken:      config.GetEnv().PrivateToken,
		GitProjectApiAddress: config.GetEnv().GitLabOperateProjectApiAddress,
		Res:                  "ok",
		Msg:                  "ok",
	}
}

func (g *GitLab) GetProjectName(HttpUrlToRepo string) *GitLab {
	projectSlice := strings.Split(HttpUrlToRepo, "/")
	projectName := projectSlice[len(projectSlice)-1:][0]
	g.ProjectName = strings.ReplaceAll(projectName, ".git", "")
	g.ProjectGitAddress = strings.ReplaceAll(HttpUrlToRepo, ".git", "") + ".git"
	log.Println("项目名称：" + g.ProjectName + ", 项目git地址：" + g.ProjectGitAddress)
	g.SetKey(HttpUrlToRepo)
	return g
}

func ProjectBranchs(c *gin.Context) {

	// gitHttpUrl := c.DefaultQuery("http_url_to_repo", "error")
	// branchName := c.DefaultQuery("branch", "error")
	gitHttpAddress := c.DefaultQuery("http_url_to_repo", "error")
	branchSlince := svc.GetBranchByRepourl(gitHttpAddress)
	// fmt.Println(gitHttpUrl, branchName)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "OK",
		"msg":  "获取成功",
		"data": branchSlince,
		// "data": 2000,
	})
}

func (g *GitLab) SetKey(appGitUrl string) *GitLab {
	re, _ := regexp.Compile(`^https://github.dm-ai.cn.*$`)
	if re.MatchString(appGitUrl) {
		g.GitPrivateToken = "5X_hTFFr4RsvUod3GPzP"
		g.GitProjectApiAddress = "https://github.dm-ai.cn/api/v4/projects"
	}
	return g
}

func (g *GitLab) GetProjectId() *GitLab {
	//
	client := resty.New()
	r, e := client.R().SetQueryParams(map[string]string{
		"private_token": g.GitPrivateToken,
		"search":        g.ProjectName,
		"simple":        "true",
		"per_page":      "1000",
	}).Get(g.GitProjectApiAddress)
	if e != nil {
		log.Error(e)
		g.Res = "fail"
		g.Msg = "查询项目id失败！"
		return g
	}

	if r.StatusCode() != 200 {
		g.Res = "fail"
		g.Msg = "查询git获得项目的id的时候，返回的http状态码不为200"
		return g
	}

	e = json.Unmarshal(r.Body(), &g.GitGetIdResponse)

	//log.Println(g.GitGetIdResponse)

	for _, v := range g.GitGetIdResponse {
		HttpUrlToRepo, _ := v["http_url_to_repo"].(string)
		if HttpUrlToRepo == g.ProjectGitAddress {
			ID, ok := v["id"].(float64)
			if !ok {
				log.Error(ok)
			}
			log.Println(ID)
			log.Println(v["id"])
			g.Id = tools.Float64ToInt64(ID)
			IdString := strconv.FormatInt(g.Id, 10)
			log.Println("项目的id：" + IdString)
			break
		}
	}

	return g
}

func (g *GitLab) GetProjectBranchs() *GitLab {

	projectBranches := make([]string, 0)

	client := resty.New()
	IdString := strconv.FormatInt(g.Id, 10)
	//IdString := strconv.FormatFloat(g.Id,'g',1,64)
	getBranchesUrl := g.GitProjectApiAddress + "/" + IdString + "/repository/branches"
	log.Println(getBranchesUrl)
	r, e := client.R().SetQueryParams(map[string]string{
		"simple":   "true",
		"per_page": "1000",
	}).SetHeader("PRIVATE-TOKEN", g.GitPrivateToken).Get(getBranchesUrl)

	if e != nil {
		log.Error(e)
		g.Res = "fail"
		g.Msg = "查询项目分支失败"
		return g
	}

	if r.StatusCode() != 200 {
		g.Res = "fail"
		g.Msg = "查询项目的分支，返回的http状态吗不为200"
		return g
	}

	e = json.Unmarshal(r.Body(), &g.GitGetBranchesResponse)

	for _, v := range g.GitGetBranchesResponse {
		branchesName, _ := v["name"].(string)
		projectBranches = append(projectBranches, branchesName)
	}
	g.ProjectBranches = projectBranches
	log.Println(projectBranches)
	return g
}
