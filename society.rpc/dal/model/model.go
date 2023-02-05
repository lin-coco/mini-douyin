package model

import "gorm.io/gorm"

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
	MsgContent string `gorm:size:256`
	FromUserId uint
	ToUserId   uint
}
