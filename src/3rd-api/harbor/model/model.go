package model

import (
	"app-deploy-platform/3rd-api/harbor/config"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type TagInfo []struct {
	Name   string `json:"name"`
	Digest string `json:"digest"`
}

type DockerHarbor struct {
	Msg          string
	Res          string
	QueryTagsURL string
	Client       *http.Client
}

func NewDockerHarbor(space, project, use, passwd string) *DockerHarbor {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				req.SetBasicAuth(use, passwd)
				return nil, nil
			},
		},
	}
	queryUrl := config.GetEnv().DockerHarborUrl + "/api/repositories/" + space + "/" + project + "/tags"
	return &DockerHarbor{
		Msg:          "OK",
		Res:          "OK",
		QueryTagsURL: queryUrl,
		Client:       client,
	}
}

func (d *DockerHarbor) GetDigest(tag string) string {
	var msg string
	resp, err := d.Client.Get(d.QueryTagsURL)
	if err != nil {

		msg = "请求docker harbor查询错误，查询url: " + d.QueryTagsURL + fmt.Sprintf("%d", resp.StatusCode)
		log.Error(msg)
		d.Res = "fail"
		d.Msg = msg
		return "500 error"
	}

	if resp.StatusCode != 200 {
		msg = "请求docker harbor查询错误，查询url: " + d.QueryTagsURL + " ,返回的状态码不为200。"
		log.Error(errors.New(msg))
		d.Res = "fail"
		d.Msg = msg
		return "500 error"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var tagInfos TagInfo
	err = json.Unmarshal(body, &tagInfos)
	if err != nil {
		msg = "对请求的结果进行反序列化失败："
		log.Error(msg, err)
		d.Res = "fail"
		d.Msg = msg
		return "500 error"
	}

	if len(tagInfos) == 0 {
		msg = "不存在任何tag"
		log.Error(errors.New(msg))
		d.Res = "fail"
		d.Msg = msg
		return "500 error"
	}

	for _, v := range tagInfos {
		if v.Name == tag {
			return v.Digest
		}
	}
	return "500 error"
}
