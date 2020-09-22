package model

import (
	"time"
)

type Space struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" time_local:"1"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
}

func (Space) TableName() string {
	return "space"
}

type GetSpace struct {
	Name string `form:"name"`
	Page int64  `form:"page"`
	Size int64  `form:"size"`
}

func NewSpace() *Space {
	s := &Space{}
	if !Model.HasTable(s.TableName()) {
		Model.CreateTable(s)
	}
	return s
}
