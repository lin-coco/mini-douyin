package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"size:256"`
	Password string `gorm:"size:256"`
}

type Video struct {
	gorm.Model
	UserId   uint
	PlayUrl  string `gorm:"size:1000"`
	CoverUrl string `gorm:"size:1000"`
	Title    string `gorm:"size:50"`
}
