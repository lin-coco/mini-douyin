package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255"`
	Password string `gorm:"size:255"`
}

type Video struct {
	gorm.Model
	UserId   uint
	PlayUrl  string `gorm:"size:1000"`
	CoverUrl string `gorm:"size:1000"`
	Title    string `gorm:"size:50"`
}

type UserFavourite struct {
	gorm.Model
	UserId  uint
	VideoId uint
	Status  uint8 `gorm:"default:1"` //点赞 状态为1 取消赞状态为0
}

type Comment struct {
	gorm.Model
	Content    string `gorm:"size:255"`
	FromUserId uint
	ToVideoId  uint
}

type Relation struct {
	gorm.Model
	FromUserId uint
	ToUserId   uint
	// 查询关注列表 即 select to_user_id from relation where from_user_id = ?
	// 查询粉丝列表 即 select from_user_id from relation where to_user_id = ?
	RelType uint8 `gorm:"default:1"` //1为有效 0为无效
}
type MessageChat struct {
	gorm.Model
	// FromUserId 给 ToUserId 发送的 MsgContent
	MsgContent string `gorm:"size:255"`
	FromUserId uint
	ToUserId   uint
}
