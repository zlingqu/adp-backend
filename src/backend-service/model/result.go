package model

import "time"

type Result struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	Name      string    `json:"name" gorm:"type:varchar(80)"`
	DeployEnv string    `json:"deploy_env" gorm:"type:varchar(80)"`
	Version   string    `json:"version" gorm:"type:varchar(160)"`
}

func (Result) TableName() string {
	return "result"
}

func NewResult() *Result {
	r := &Result{}
	if !Model.HasTable(r.TableName()) {
		Model.CreateTable(r)
	}
	return r
}
