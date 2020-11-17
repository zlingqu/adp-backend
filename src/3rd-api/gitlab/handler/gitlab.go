package handler

import (
	svc "app-deploy-platform/3rd-api/gitlab/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCommitID(c *gin.Context) {

	gitHttpUrl := c.DefaultQuery("http_url_to_repo", "error")
	branchName := c.DefaultQuery("branch", "error")
	gitHttpUrl = svc.GitlabUrlCheck(gitHttpUrl) //url格式检查和转换
	commitID := svc.GetCommitIDByRepourlAndBranch(gitHttpUrl, branchName)
	if commitID == "error" {
		c.JSON(http.StatusForbidden, gin.H{
			"res":  "error",
			"msg":  "找不到这样的repo或者分支名不正确,请检查",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "OK",
		"msg":  "获取成功",
		"data": commitID,
	})

}

func GetTags(c *gin.Context) {

	gitHttpUrl := c.DefaultQuery("http_url_to_repo", "error")
	gitHttpUrl = svc.GitlabUrlCheck(gitHttpUrl) //url格式检查和转换
	tagSlince := svc.GetTagsByRepourl(gitHttpUrl)
	if tagSlince == nil {
		c.JSON(http.StatusForbidden, gin.H{
			"res":  "error",
			"msg":  "找不到这样的repo,请检查",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "OK",
		"msg":  "获取成功",
		"data": tagSlince,
	})

}

func GetBranchs(c *gin.Context) {

	gitHttpUrl := c.DefaultQuery("http_url_to_repo", "error")
	gitHttpUrl = svc.GitlabUrlCheck(gitHttpUrl) //url格式检查和转换
	branchSlince := svc.GetBranchByRepourl(gitHttpUrl)
	if branchSlince == nil {
		c.JSON(http.StatusForbidden, gin.H{
			"res":  "error",
			"msg":  "找不到这样的repo,请检查",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"res":  "OK",
		"msg":  "获取成功",
		"data": branchSlince,
	})

}

func ProjectInfo(c *gin.Context) {
	gitHttpUrl := c.DefaultQuery("http_url_to_repo", "error")
	gitHttpUrl = svc.GitlabUrlCheck(gitHttpUrl) //url格式检查和转换
	id := svc.GetIDByRepourl(gitHttpUrl)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"res":  "error",
			"msg":  "找不到这样的repo,请检查",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"res":  "OK",
		"data": id,
	})

}
