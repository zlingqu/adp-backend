package model

import (
	"app-deploy-platform/backend-service/config"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"time"
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
	if !Model.HasTable(u.TableName()) {
		Model.CreateTable(u)
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
	client := resty.New()
	misLdapServiceUrl := config.GetEnv().MisLdapServiceUrl
	r, e := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(m.Req).Post(misLdapServiceUrl)

	// 处理请求问题和错误
	if e != nil {
		log.Error(e)
		m.Res = "fail"
		m.Msg = "请求mis-ldap-service异常。"
		return m
	}

	if r.StatusCode() != 200 {
		m.Res = "fail"
		m.Msg = "请求mis-ldap-service结果状态码，不等于200。"
		return m
	}

	e = json.Unmarshal(r.Body(), &m.Rep)

	if e != nil {
		log.Error(e)
		m.Res = "fail"
		m.Msg = "对mis-ldap-service的请求结果进行json反序列化失败。"
		return m
	}
	return m
}
