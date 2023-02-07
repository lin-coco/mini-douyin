package handler

import "api.service/biz/model/api/douyin/core"

type RVideo struct {
	Id            int64  // 视频唯一标识 18
	UserId        int64  // 视频作者id
	PlayUrl       string // 视频播放地址
	CoverUrl      string // 视频封面地址
	FavoriteCount int64  // 视频的点赞总数
	CommentCount  int64  // 视频的评论总数
	Title         string // 视频标题
}

type RUser struct {
	Id            int64  // 用户id
	Name          string // 用户名称
	FollowCount   int64  // 关注总数
	FollowerCount int64  // 粉丝总数
}

func VideoToR(videos []*core.Video) ([]RVideo, []RUser) {
	rVideos := make([]RVideo, 0, len(videos))
	rUsers := make([]RUser, 0, len(videos))
	for _, video := range videos {
		rVideos = append(rVideos, RVideo{
			Id:            video.Id,
			UserId:        video.Author.Id,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
		})
		rUsers = append(rUsers, RUser{
			Id:            video.Author.Id,
			Name:          video.Author.Name,
			FollowCount:   *video.Author.FollowCount,
			FollowerCount: *video.Author.FollowerCount,
		})
	}
	return rVideos, rUsers
}
func RToVideo(rVideos []RVideo) []*core.Video {
	return nil
}

func UserToR(users []*core.User) []RUser {
	//rUsers := make([]RVideo, 0, len(users))
	//for _, video := range users {
	//	rUsers = append(rUsers, RVideo{})
	//}
	//return rUsers
	return nil
}

func RToUser(rUsers []RUser) []*core.User {
	return nil
}
