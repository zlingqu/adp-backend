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
	ID         int    `json:"id"`
	Httpurl    string `json:"http_url_to_repo"`
	BranchName string `json:"name"`
}

type ProBranchStr struct {
	BranchName string `json:"name"`
}

type ProCommitStr struct {
	ID string `json:"id"`
}

func GetBranchByRepourl(httpurl string) []string {
	arr := strings.Split(httpurl, "/")
	urlshort := arr[0] + "//" + arr[2]
	proID := getIDByRepourl(httpurl)
	proIDStr := strconv.Itoa(proID)
	var branchsMap []ProBranchStr
	fillDataByGitlab(urlshort+"/api/v4/projects/"+proIDStr+"/repository/branches", &branchsMap)
	var branchsSlicer []string
	for _, v := range branchsMap {
		// branchsSlicer.Append
		branchsSlicer = append(branchsSlicer, v.BranchName)

	}
	return branchsSlicer
}

func GetCommitIDByRepourlAndBranch(httpurl, branchName string) string {
	arr := strings.Split(httpurl, "/")
	urlshort := arr[0] + "//" + arr[2]
	proID := getIDByRepourl(httpurl)
	proIDStr := strconv.Itoa(proID)

	var commits ProCommitStr
	fillDataByGitlab(urlshort+"/api/v4/projects/"+proIDStr+"/repository/commits/"+branchName, &commits)

	return commits.ID
}

func getIDByRepourl(httpurl string) int {
	arr := strings.Split(httpurl, "/")
	urlshort := arr[0] + "//" + arr[2]
	proSlince := strings.Split(arr[len(arr)-1], ".")
	proName := proSlince[0]

	var ids []ProStr
	fillDataByGitlab(urlshort+"/api/v4/projects/?search="+proName, &ids)

	for _, v := range ids {
		if v.Httpurl == httpurl {
			return v.ID
		}

	}
	return 0
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
