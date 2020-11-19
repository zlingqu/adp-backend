package service

import (
	cfg "app-deploy-platform/3rd-api/gitlab/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type ProStr struct {
	ID      int    `json:"id"`
	Httpurl string `json:"http_url_to_repo"`
	ProName string `json:"name"`
}

type ProBranchStr struct {
	BranchName string `json:"name"`
}

type ProCommitStr struct {
	ID string `json:"id"`
}

type ProTagStr struct {
	TagName string `json:"name"`
}

func GetBranchByRepourl(httpurl string) []string {
	arr := strings.Split(httpurl, "/")
	urlshort := arr[0] + "//" + arr[2]
	proID := GetIDByRepourl(httpurl)
	if proID == -1 {
		return nil
	}
	proIDStr := strconv.Itoa(proID)
	var branchsMap []ProBranchStr
	fillDataByGitlab(urlshort+"/api/v4/projects/"+proIDStr+"/repository/branches", &branchsMap)
	var branchsSlicer []string
	for _, v := range branchsMap {
		branchsSlicer = append(branchsSlicer, v.BranchName)

	}
	return branchsSlicer
}

func GetTagsByRepourl(httpurl string) []string {
	arr := strings.Split(httpurl, "/")
	urlshort := arr[0] + "//" + arr[2]
	proID := GetIDByRepourl(httpurl)
	if proID == -1 {
		return nil
	}
	proIDStr := strconv.Itoa(proID)
	var tagMap []ProTagStr

	fillDataByGitlab(urlshort+"/api/v4/projects/"+proIDStr+"/repository/tags", &tagMap)
	if len(tagMap) == 0 {
		return []string{"null,没有tag"} //没有tag
	}
	var tagSlicer []string
	for _, v := range tagMap {
		tagSlicer = append(tagSlicer, v.TagName)

	}

	return tagSlicer
}

func GetCommitIDByRepourlAndBranch(httpurl, branchName string) string {
	arr := strings.Split(httpurl, "/")
	urlshort := arr[0] + "//" + arr[2]
	proID := GetIDByRepourl(httpurl)
	if proID == -1 {
		return "error"
	}
	proIDStr := strconv.Itoa(proID)

	var commits ProCommitStr
	fillDataByGitlab(urlshort+"/api/v4/projects/"+proIDStr+"/repository/commits/"+branchName, &commits)
	if commits.ID == "" {
		return "error"
	}

	return commits.ID
}

func GetIDByRepourl(httpurl string) int {
	arr := strings.Split(httpurl, "/")
	urlshort := arr[0] + "//" + arr[2]
	proSlicer := strings.Split(arr[len(arr)-1], ".")
	proName := proSlicer[0]

	var ids []ProStr
	fillDataByGitlab(urlshort+"/api/v4/projects/?search="+proName+"&simple=true&per_page=300", &ids)

	for _, v := range ids {
		if v.Httpurl == httpurl {
			return v.ID
		}

	}
	return -1 //-1表示不存在对应的项目
}

func GitlabUrlCheck(url string) string {
	if !strings.HasSuffix(url, ".git") { //如果url写错，没有以.git结尾，将其加上
		url = url + ".git"
	}
	if strings.HasPrefix(url, "git@") { //如果填写的是ssh协议的，转换成https
		arr := strings.Split(url, ":")
		sshShortUrl := arr[0]
		pubUrl := arr[1]
		domain := strings.Split(sshShortUrl, "@")[1]
		url = "https://" + domain + "/" + pubUrl
	}
	return url
}

func fillDataByGitlab(url string, ptr interface{}) {
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)

	arr := strings.Split(url, "/")
	urlshort := arr[0] + "//" + arr[2]
	re, _ := regexp.Compile(`^https://github.dm-ai.cn.*$`)
	if re.MatchString(urlshort) {
		reqest.Header.Add("PRIVATE-TOKEN", "5X_hTFFr4RsvUod3GPzP")
	} else {
		reqest.Header.Add("PRIVATE-TOKEN", cfg.GetEnv().PrivateToken)
	}

	resp, err := client.Do(reqest)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, ptr)
}
