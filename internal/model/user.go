package model

import "time"

// users
type User struct {
	Id         int           `json:"id" gorm:"column:id"`
	Name       string        `json:"name" gorm:"column:name"`
	Cet4       int           `json:"cet4" gorm:"column:cet4"`
	Cet6       int           `json:"cet6" gorm:"column:cet6"`
	CreateAt   time.Duration `json:"createAt"`
	UpdateAt   time.Duration `json:"updateAt"`
	DeleteFlag int           `json:"deleteFlag" gorm:"delete_flag"`
}
