package model

import (
	"gorm.io/gorm"
)

type UserFavourite struct {
	gorm.Model
	UserId  uint
	VideoId uint
	Status  uint8 `gorm:"default:1"` //点赞 状态为1 取消赞状态为0
}

type Comment struct {
	gorm.Model
	Content    string `gorm:"size:256"`
	FromUserId uint
	ToVideoId  uint
}
