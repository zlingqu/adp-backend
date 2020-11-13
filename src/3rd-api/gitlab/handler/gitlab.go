package handler

import (
	svc "app-deploy-platform/3rd-api/gitlab/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func ProjectBranchs(c *gin.Context) {

	gitHttpAddress := c.DefaultQuery("http_url_to_repo", "error")
	branchSlince := svc.GetBranchByRepourl(gitHttpAddress)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "OK",
		"msg":  "获取成功",
		"data": branchSlince,
	})
}
