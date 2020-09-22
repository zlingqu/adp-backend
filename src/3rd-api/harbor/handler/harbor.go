package handler

import (
	"app-deploy-platform/3rd-api/harbor/config"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"net/http"
)

func DockerImageSha256(c *gin.Context) {

	space := c.DefaultQuery("space", "")
	project := c.DefaultQuery("project", "")
	tag := c.DefaultQuery("tag", "")
	//queryUrl := conf.GetEnv().DockerHarborUrl + "/api/repositories/" + space + "/" + project + "/tags?detail=1"

	dh := NewDockerHarbor()
	dh.Init(space, project).QueryTags(tag)

	repStatsCode := http.StatusOK
	if dh.Res != "ok" {
		repStatsCode = http.StatusInternalServerError
	}

	c.JSON(repStatsCode, gin.H{
		"code": 0,
		"res":  dh.Res,
		"msg":  dh.Msg,
		"data": dh.QueryDigest,
	})

}

func NewDockerHarbor() *DockerHarbor {
	return &DockerHarbor{}
}

type DockerHarbor struct {
	QueryDigest  string
	Msg          string
	Res          string
	QueryTagsUrl string
	TagInfo      []struct {
		Name   string `json:"name"`
		Digest string `json:"digest"`
	}
}

func (d *DockerHarbor) Init(space string, project string) *DockerHarbor {
	d.QueryTagsUrl = config.GetEnv().DockerHarborUrl + "/api/repositories/" + space + "/" + project + "/tags?detail=1"
	d.Msg = "ok"
	d.Res = "ok"
	d.QueryDigest = ""
	return d
}

func (d *DockerHarbor) QueryTags(tag string) *DockerHarbor {
	client := resty.New()
	r, e := client.R().SetHeader("Accept", "application/json").SetBasicAuth(config.GetEnv().DockerHarborUser, config.GetEnv().DockerHarborPassword).Get(d.QueryTagsUrl)
	var msg string

	if e != nil {
		msg = "请求docker harbor查询错误，查询url: " + d.QueryTagsUrl
		log.Error(msg)
		d.Res = "fail"
		d.Msg = msg
		return d
	}

	if r.StatusCode() != 200 {
		msg = "请求docker harbor查询错误，查询url: " + d.QueryTagsUrl + " ,返回的状态码不为200。"
		log.Error(errors.New(msg))
		d.Res = "fail"
		d.Msg = msg
		return d
	}

	e = json.Unmarshal(r.Body(), &d.TagInfo)
	if e != nil {
		msg = "对请求的结果进行反序列化失败："
		log.Error(msg, e)
		d.Res = "fail"
		d.Msg = msg
		return d
	}

	log.Info("Docker harbor 返回的结果信息为：", d.TagInfo)

	// 检查返回的结果是否为空
	if len(d.TagInfo) == 0 {
		msg = "查询的tag为空。"
		log.Error(errors.New(msg))
		d.Res = "fail"
		d.Msg = msg
		return d
	}

	// 轮训结果，找到tag对应的sha值
	for _, v := range d.TagInfo {
		if v.Name == tag {
			d.QueryDigest = v.Digest
			return d
		}
	}

	return d
}
