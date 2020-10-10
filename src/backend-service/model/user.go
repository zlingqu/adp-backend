package model

import (
	"app-deploy-platform/backend-service/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type UserInfo struct {
	OwnerEnglishName string `json:"owner_english_name" gorm:"type:varchar(80);primary_key"`
	OwnerChinaName   string `json:"owner_china_name" gorm:"type:varchar(80)"`
}

type User struct {
	//ID               uint      `json:"id" gorm:"primary_key"`
	CreatedAt        time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	UpdatedAt        time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	OwnerEnglishName string    `json:"owner_english_name" gorm:"type:varchar(80);primary_key"`
	OwnerChinaName   string    `json:"owner_china_name" gorm:"type:varchar(80)"`
	Status           string    `json:"status" gorm:"type:varchar(80)"`
}

func (User) TableName() string {
	return "user"
}

func NewUser() *User {
	u := &User{}
	// if !Model.HasTable(u.TableName()) {
	// 	Model.CreateTable(u)
	// }
	if Model.HasTable(u.TableName()) { //判断表是否存在
		Model.AutoMigrate(u) //存在就自动适配表，也就说原先没字段的就增加字段
	} else {
		Model.CreateTable(u) //不存在就创建新表
	}
	return u
}

type LdapUserInfo struct {
	Cn          string `json:"cn"`
	DisplayName string `json:"displayName"`
}

type MisLdapServiceReq struct {
	Type      string `json:"type"`
	Condition struct {
		L      string `json:"l"`
		IsShow string `json:"isShow"`
	} `json:"condition"`
	Attributes []string `json:"attributes"`
}

type MisLdapServiceRep struct {
	Code  int            `json:"code"`
	Data  []LdapUserInfo `json:"data"`
	error interface{}    `json:"error"`
}

type MisLdapService struct {
	Req MisLdapServiceReq
	Rep MisLdapServiceRep
	Res string
	Msg string
}

/*
mis的ladp接口需要的接口格式
{
    "type":"people",
    "condition":{
        "l":"中国-广州",
        "isShow":"1"
    },
    "attributes":[
        "cn",
        "displayName"
    ]
}
*/

func NewMisLdapService() *MisLdapService {
	return &MisLdapService{
		Req: MisLdapServiceReq{
			Type: "people",
			Condition: struct {
				L      string `json:"l"`
				IsShow string `json:"isShow"`
			}{
				L:      "中国-广州",
				IsShow: "1",
			},
			Attributes: []string{"cn", "displayName"},
		},
		Rep: MisLdapServiceRep{},
		Res: "ok",
		Msg: "Search mis-ldap-service succeed",
	}
}

func (m *MisLdapService) Request() *MisLdapService {

	misLdapServiceURL := config.GetEnv().MisLdapServiceUrl

	jsonStr, err := json.Marshal(m.Req)
	if err != nil {
		fmt.Println("json转换错误")
	}
	req, _ := http.NewRequest("POST", misLdapServiceURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Error("请求mis-ldap-service,错误信息是", err)
		m.Res = "fail"
		m.Msg = "请求mis-ldap-service异常。"
		return m
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &m.Rep)
	statuscode := resp.StatusCode
	if statuscode != 200 {
		m.Res = "fail"
		m.Msg = "请求mis-ldap-service结果状态码，不等于200。"
		return m
	}

	if err != nil {
		m.Res = "fail"
		m.Msg = "对mis-ldap-service的请求结果进行json反序列化失败。"
		return m
	}

	return m
}
