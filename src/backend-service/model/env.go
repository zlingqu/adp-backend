package model

import (
	"time"
)

type Env struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
}

func (Env) TableName() string {
	return "env"
}

func NewEnv() *Env {
	e := &Env{}
	// if !Model.HasTable(e.TableName()) {
	// 	Model.CreateTable(e)
	// }
	if Model.HasTable(e.TableName()) { //判断表是否存在
		Model.AutoMigrate(e) //存在就自动适配表，也就说原先没字段的就增加字段
	} else {
		Model.CreateTable(e) //不存在就创建新表
	}
	return e
}

type GetEnv struct {
	Name string `form:"name"`
	Page int64  `form:"page"`
	Size int64  `form:"size"`
}

type GetEnvByID struct {
	ID string `uri:"id" binding:"required"`
}

type PostEnvs struct {
	Ids []string `json:"ids"`
}
