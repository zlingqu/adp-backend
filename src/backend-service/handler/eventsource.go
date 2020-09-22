package handler

import (
	"app-deploy-platform/backend-service/config"
	"encoding/json"
	"errors"
	"gopkg.in/resty.v1"
	"time"
)

type Result struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	Name      string    `json:"name" gorm:"type:varchar(80)"`
	DeployEnv string    `json:"deploy_env" gorm:"type:varchar(80)"`
	Version   string    `json:"version" gorm:"type:varchar(160)"`
}

type ResponseResult struct {
	Code int    `json:"code"`
	Data Result `json:"data"`
	Msg  string `json:"msg"`
	Res  string `json:"res"`
}

func SearchAdpResultInfo(deployEnv string, appName string) (Result, error) {
	// init
	var getResult ResponseResult

	serviceAdpBuildResultUrl := config.GetEnv().ServiceAdpBuildResultUrl
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"deployEnv": deployEnv,
			"name":      appName,
		}).Get(serviceAdpBuildResultUrl)

	if err != nil {
		return Result{}, err
	}

	if resp.StatusCode() != 200 {
		return Result{}, errors.New("请求serviceAdpBuildResult http的返回状态码不为200。")
	}

	err = json.Unmarshal(resp.Body(), &getResult)

	if err != nil {
		return Result{}, errors.New("请求serviceAdpBuildResult成功后，对结果就行反序列化失败。")
	}

	if getResult.Res != "ok" {
		return Result{}, errors.New(getResult.Res)
	}

	if getResult.Data.ID == 0 {
		return Result{}, errors.New("service-adp-build-result的id返回为0")
	}

	return getResult.Data, nil
}
