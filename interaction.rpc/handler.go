package main

import (
	"basics.rpc/kitex_gen/douyin/core"
	"context"
	"errors"
	"interaction.rpc/dal/model"
	"interaction.rpc/dal/query"
	first "interaction.rpc/kitex_gen/douyin/extra/first"
	"log"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct{}

// AddVideoFavorite implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) AddVideoFavorite(ctx context.Context, req *first.AddVideoFavoriteRequest) (resp *first.AddVideoFavoriteResponse, err error) {
	userId := req.UserId
	videoId := req.VideoId
	userFavourite, _ := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).First()
	if userFavourite == nil {
		err = Q.UserFavourite.Create(&model.UserFavourite{UserId: uint(userId), VideoId: uint(videoId), Status: 1})
		if err != nil {
			log.Printf("create userFavorite failed err:%v", err)
			return nil, err
		}
	} else {
		_, err := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).Update(Q.UserFavourite.Status, 1)
		if err != nil {
			log.Printf("update userFavorite failed err:%v", err)
		}
	}
	return &first.AddVideoFavoriteResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// CancelVideoFavorite implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CancelVideoFavorite(ctx context.Context, req *first.CancelVideoFavoriteRequest) (resp *first.CancelVideoFavoriteResponse, err error) {
	userId := req.UserId
	videoId := req.VideoId
	userFavourite, _ := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).First()
	if userFavourite != nil {
		_, err := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).Update(Q.UserFavourite.Status, 0)
		if err != nil {
			log.Printf("cancel favorite failed userId:%d err:%s", userId, err)
			return nil, err
		}
		return &first.CancelVideoFavoriteResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		}, nil
	}
	log.Printf("user has canceled")
	return &first.CancelVideoFavoriteResponse{
		StatusCode: 0,
		StatusMsg:  "success"}, nil
}

// GetFavoriteList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetFavoriteList(ctx context.Context, req *first.GetFavoriteListRequest) (resp *first.GetFavoriteListResponse, err error) {
	userId := req.UserId
	userFavourites, err := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.Status.Eq(1)).Find()
	if err != nil {
		log.Printf("query userFavourites failed userId:%d err:%s", userId, err)
		return nil, err
	}
	//videoList := make([]*first.Video, 0, len(userFavourites))
	videoIdList := make([]int64, 0, len(userFavourites))
	for _, userFavourite := range userFavourites {
		videoId := userFavourite.VideoId
		videoIdList = append(videoIdList, int64(videoId))
	}
	if len(videoIdList) == 0 {
		return &first.GetFavoriteListResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			VideoList:  make([]*first.Video, 0, 0),
		}, nil
	}
	res, err := BasicsService.GetVideoListByIds(ctx, &core.GetVideoListByIdsRequest{VideoIdList: videoIdList})
	if err != nil {
		log.Printf("BasicsService failed err:%v", err)
		return nil, err
	}
	if res.StatusCode != 0 {
		log.Println(res.StatusMsg)
		return nil, errors.New(res.StatusMsg)
	}
	videos := res.VideoList

	videoList := make([]*first.Video, 0, len(videos))
	for _, video := range videos {
		//query favorite Count
		favoriteCount, err := Q.UserFavourite.Where(query.UserFavourite.VideoId.Eq(uint(video.Id)), query.UserFavourite.Status.Eq(1)).Count()
		if err != nil {
			log.Printf("query favorite count failed err:%v", err)
			return nil, err
		}
		//query comment count
		commentCount, err := Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(video.Id))).Count()
		if err != nil {
			log.Printf("query comment count failed err:%v", err)
			return nil, err
		}

		v := &first.Video{
			Id:            video.Id,
			Author:        &first.User{Id: video.User.Id, Name: video.User.Name},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			Title:         video.Title,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    true,
		}
		videoList = append(videoList, v)
	}

	return &first.GetFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}, nil
}

// AddComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) AddComment(ctx context.Context, req *first.AddCommentRequest) (resp *first.AddCommentResponse, err error) {
	userId := req.UserId
	videoId := req.VideoId
	commentText := req.CommentText

	res, err := BasicsService.GetUserInfoById(ctx, &core.GetUserRequest{UserId: userId})
	if err != nil {
		log.Printf("BasicsService run failed err:%v", err)
	}

	comment := &model.Comment{
		FromUserId: uint(userId),
		ToVideoId:  uint(videoId),
		Content:    commentText,
	}
	err = Q.Comment.Create(comment)
	if err != nil {
		log.Printf("create comment failed err:%v", err)
		return nil, err
	}
	return &first.AddCommentResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Comment: &first.Comment{
			Id:         int64(comment.ID),
			User:       &first.User{Id: userId, Name: res.Name},
			Content:    commentText,
			CreateDate: comment.CreatedAt.Format("01-02"),
		},
	}, nil
}

// DeleteComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) DeleteComment(ctx context.Context, req *first.DeleteCommentRequest) (resp *first.DeleteCommentResponse, err error) {
	commentId := req.CommentId

	_, err = Q.Comment.Where(query.Comment.ID.Eq(uint(commentId))).Delete()
	if err != nil {
		log.Printf("delete comment failed commentId:%d err:%v", commentId, err)
		return nil, err
	}

	return &first.DeleteCommentResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// CommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentList(ctx context.Context, req *first.CommentListRequest) (resp *first.CommentListResponse, err error) {
	videoId := req.VideoId
	comments, err := Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(videoId))).Order().Find()
	if err != nil {
		log.Printf("query comments failed videoId:%d err:%v", videoId, comments)
		return nil, err
	}
	commentList := make([]*first.Comment, 0, len(comments))
	for _, comment := range comments {
		userId := comment.FromUserId
		res, err := BasicsService.GetUserInfoById(ctx, &core.GetUserRequest{UserId: int64(userId)})
		if err != nil {
			log.Printf("BasicsService query failed err:%v", err)
		}
		commentList = append(commentList, &first.Comment{
			Id:         int64(comment.ID),
			User:       &first.User{Id: int64(userId), Name: res.Name},
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		})
	}
	return &first.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		CommentList: commentList,
	}, nil
}

// GetVideoFavoriteCount implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetVideoFavoriteCount(ctx context.Context, req *first.GetVideoFavoriteCountRequest) (resp *first.GetVideoFavoriteCountResponse, err error) {
	videoId := req.VideoId
	count, err := Q.UserFavourite.Where(query.UserFavourite.VideoId.Eq(uint(videoId)), query.UserFavourite.Status.Eq(1)).Count()
	if err != nil {
		log.Printf("video favorite count query failed err:%v", err)
		return nil, err
	}
	return &first.GetVideoFavoriteCountResponse{
		StatusCode:    0,
		StatusMsg:     "success",
		FavoriteCount: count,
	}, nil
}

// GetVideoCommentCount implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetVideoCommentCount(ctx context.Context, req *first.GetVideoCommentCountRequest) (resp *first.GetVideoCommentCountResponse, err error) {
	videoId := req.VideoId
	count, err := Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(videoId))).Count()
	if err != nil {
		log.Printf("video comment count query failed err:%v", err)
		return nil, err
	}
	return &first.GetVideoCommentCountResponse{
		StatusCode:   0,
		StatusMsg:    "success",
		CommentCount: count,
	}, nil
}

// IsFavorite implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) IsFavorite(ctx context.Context, req *first.IsFavoriteRequest) (resp *first.IsFavoriteResponse, err error) {
	userId := req.UserId
	videoId := req.VideoId
	favourite, _ := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId)), query.UserFavourite.Status.Eq(1)).First()
	if favourite == nil {
		return &first.IsFavoriteResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			IsFavorite: false,
		}, nil
	}
	return &first.IsFavoriteResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		IsFavorite: true,
	}, nil
}

// GetCommentById implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetCommentById(ctx context.Context, req *first.GetCommentByIdRequest) (resp *first.GetCommentByIdResponse, err error) {
	commentId := req.CommentId
	comment, err := Q.Comment.Where(query.Comment.ID.Eq(uint(commentId))).First()
	if comment == nil {
		log.Printf("commentId is not existed")
		return nil, err
	}

	return &first.GetCommentByIdResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		CommentId:  commentId,
		UserId:     int64(comment.FromUserId),
		Content:    comment.Content,
		VideoId:    int64(comment.ToVideoId),
	}, nil
}
