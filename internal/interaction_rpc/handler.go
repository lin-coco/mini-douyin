package interaction_rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"mini-douyin/internal/pkg/dal/model"
	"mini-douyin/internal/pkg/dal/query"
	"mini-douyin/internal/pkg/kitex_gen/douyin/basics"
	"mini-douyin/internal/pkg/kitex_gen/douyin/interaction"
	"mini-douyin/pkg/cache"
	"time"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct{}

// AddVideoFavorite implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) AddVideoFavorite(ctx context.Context, req *interaction.AddVideoFavoriteRequest) (resp *interaction.AddVideoFavoriteResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	videoId := req.VideoId
	userFavourite, _ := query.Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).First()
	if userFavourite == nil {
		err = query.Q.UserFavourite.Create(&model.UserFavourite{UserId: uint(userId), VideoId: uint(videoId), Status: 1})
		if err != nil {
			klog.Infof("create userFavorite failed err:%v", err)
			return nil, err
		}
	} else {
		_, err := query.Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).Update(query.Q.UserFavourite.Status, 1)
		if err != nil {
			klog.Infof("update userFavorite failed err:%v", err)
		}
	}
	//cache
	result, _ := cache.Client.Exists(ctx, fmt.Sprintf("interaction.rpc:videoFavoriteCount:%d", videoId)).Result()
	if result > 0 {
		cache.Client.Incr(ctx, fmt.Sprintf("interaction.rpc:videoFavoriteCount:%d", videoId))
	}
	return &interaction.AddVideoFavoriteResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// CancelVideoFavorite implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CancelVideoFavorite(ctx context.Context, req *interaction.CancelVideoFavoriteRequest) (resp *interaction.CancelVideoFavoriteResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	videoId := req.VideoId
	userFavourite, _ := query.Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).First()
	if userFavourite != nil {
		_, err := query.Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).Update(query.Q.UserFavourite.Status, 0)
		if err != nil {
			klog.Infof("cancel favorite failed userId:%d err:%s", userId, err)
			return nil, err
		}
		return &interaction.CancelVideoFavoriteResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		}, nil
	}
	klog.Infof("user has canceled")
	//cache
	result, _ := cache.Client.Exists(ctx, fmt.Sprintf("interaction.rpc:videoFavoriteCount:%d", videoId)).Result()
	if result > 0 {
		cache.Client.Decr(ctx, fmt.Sprintf("interaction.rpc:videoFavoriteCount:%d", videoId))
	}
	return &interaction.CancelVideoFavoriteResponse{
		StatusCode: 0,
		StatusMsg:  "success"}, nil
}

// GetFavoriteList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetFavoriteList(ctx context.Context, req *interaction.GetFavoriteListRequest) (resp *interaction.GetFavoriteListResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	userFavourites, err := query.Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.Status.Eq(1)).Find()
	if err != nil {
		klog.Infof("query userFavourites failed userId:%d err:%s", userId, err)
		return nil, err
	}
	//videoList := make([]*interaction.Video, 0, len(userFavourites))
	videoIdList := make([]int64, 0, len(userFavourites))
	for _, userFavourite := range userFavourites {
		videoId := userFavourite.VideoId
		videoIdList = append(videoIdList, int64(videoId))
	}
	if len(videoIdList) == 0 {
		return &interaction.GetFavoriteListResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			VideoList:  make([]*interaction.Video, 0, 0),
		}, nil
	}
	res, err := BasicsRpcClient.GetVideoListByIds(ctx, &basics.GetVideoListByIdsRequest{VideoIdList: videoIdList})
	if err != nil {
		klog.Infof("BasicsRpcClient failed err:%v", err)
		return nil, err
	}
	if res.StatusCode != 0 {
		klog.Infof(res.StatusMsg)
		return nil, errors.New(res.StatusMsg)
	}
	videos := res.VideoList

	videoList := make([]*interaction.Video, 0, len(videos))
	for _, video := range videos {
		//query favorite Count
		favoriteCount, err := query.Q.UserFavourite.Where(query.UserFavourite.VideoId.Eq(uint(video.Id)), query.UserFavourite.Status.Eq(1)).Count()
		if err != nil {
			klog.Infof("query favorite count failed err:%v", err)
			return nil, err
		}
		//query comment count
		commentCount, err := query.Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(video.Id))).Count()
		if err != nil {
			klog.Infof("query comment count failed err:%v", err)
			return nil, err
		}

		v := &interaction.Video{
			Id:            video.Id,
			Author:        &interaction.User{Id: video.User.Id, Name: video.User.Name},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			Title:         video.Title,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    true,
		}
		videoList = append(videoList, v)
	}

	return &interaction.GetFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}, nil
}

// AddComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) AddComment(ctx context.Context, req *interaction.AddCommentRequest) (resp *interaction.AddCommentResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	videoId := req.VideoId
	commentText := req.CommentText

	res, err := BasicsRpcClient.GetUserInfoById(ctx, &basics.GetUserRequest{UserId: userId})
	if err != nil {
		klog.Infof("BasicsRpcClient run failed err:%v", err)
	}

	comment := &model.Comment{
		FromUserId: uint(userId),
		ToVideoId:  uint(videoId),
		Content:    commentText,
	}
	err = query.Q.Comment.Create(comment)
	if err != nil {
		klog.Infof("create comment failed err:%v", err)
		return nil, err
	}
	//cache
	result, _ := cache.Client.Exists(ctx, fmt.Sprintf("interaction.rpc:videoCommentCount:%d", videoId)).Result()
	if result > 0 {
		cache.Client.Incr(ctx, fmt.Sprintf("interaction.rpc:videoCommentCount:%d", videoId))
	}
	return &interaction.AddCommentResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Comment: &interaction.Comment{
			Id:         int64(comment.ID),
			User:       &interaction.User{Id: userId, Name: res.Name},
			Content:    commentText,
			CreateDate: comment.CreatedAt.Format("01-02"),
		},
	}, nil
}

// DeleteComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) DeleteComment(ctx context.Context, req *interaction.DeleteCommentRequest) (resp *interaction.DeleteCommentResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	commentId := req.CommentId

	comment, err := query.Q.Comment.Where(query.Comment.ID.Eq(uint(commentId))).Select(query.Comment.ToVideoId).First()
	if err != nil {
		klog.Infof("delete comment failed commentId:%d err:%v", commentId, err)
		return nil, err
	}
	_, err = query.Q.Comment.Where(query.Comment.ID.Eq(uint(commentId))).Delete()
	if err != nil {
		klog.Infof("delete comment failed commentId:%d err:%v", commentId, err)
		return nil, err
	}
	//cache
	result, _ := cache.Client.Exists(ctx, fmt.Sprintf("interaction.rpc:videoCommentCount:%d", comment.ToVideoId)).Result()
	if result > 0 {
		cache.Client.Decr(ctx, fmt.Sprintf("interaction.rpc:videoCommentCount:%d", comment.ToVideoId))
	}
	return &interaction.DeleteCommentResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// CommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentList(ctx context.Context, req *interaction.CommentListRequest) (resp *interaction.CommentListResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	videoId := req.VideoId
	comments, err := query.Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(videoId))).Order(query.Comment.CreatedAt).Find()
	if err != nil {
		klog.Infof("query comments failed videoId:%d err:%v", videoId, comments)
		return nil, err
	}
	commentList := make([]*interaction.Comment, 0, len(comments))
	for _, comment := range comments {
		userId := comment.FromUserId
		res, err := BasicsRpcClient.GetUserInfoById(ctx, &basics.GetUserRequest{UserId: int64(userId)})
		if err != nil {
			klog.Infof("BasicsRpcClient query failed err:%v", err)
		}
		commentList = append(commentList, &interaction.Comment{
			Id:         int64(comment.ID),
			User:       &interaction.User{Id: int64(userId), Name: res.Name},
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		})
	}
	return &interaction.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		CommentList: commentList,
	}, nil
}

// GetVideoFavoriteCount implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetVideoFavoriteCount(ctx context.Context, req *interaction.GetVideoFavoriteCountRequest) (resp *interaction.GetVideoFavoriteCountResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	videoId := req.VideoId
	count, err := cache.Client.Get(ctx, fmt.Sprintf("interaction.rpc:videoFavoriteCount:%d", videoId)).Int64()
	if err != nil {
		klog.Infof("cache cache not exist || cache server failed")
	} else {
		return &interaction.GetVideoFavoriteCountResponse{
			StatusCode:    0,
			StatusMsg:     "success",
			FavoriteCount: count,
		}, nil
	}
	count, err = query.Q.UserFavourite.Where(query.UserFavourite.VideoId.Eq(uint(videoId)), query.UserFavourite.Status.Eq(1)).Count()
	if err != nil {
		klog.Infof("video favorite count query failed err:%v", err)
		return nil, err
	}
	err = cache.Client.Set(ctx, fmt.Sprintf("interaction.rpc:videoFavoriteCount:%d", videoId), count, time.Hour).Err()
	if err != nil {
		klog.Infof("cache failed cache interaction.rpc:videoFavoriteCount")
	}
	return &interaction.GetVideoFavoriteCountResponse{
		StatusCode:    0,
		StatusMsg:     "success",
		FavoriteCount: count,
	}, nil
}

// GetVideoCommentCount implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetVideoCommentCount(ctx context.Context, req *interaction.GetVideoCommentCountRequest) (resp *interaction.GetVideoCommentCountResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	videoId := req.VideoId
	count, err := cache.Client.Get(ctx, fmt.Sprintf("interaction.rpc:videoCommentCount:%d", videoId)).Int64()
	if err != nil {
		klog.Infof("cache cache not exist || cache server failed")
	} else {
		return &interaction.GetVideoCommentCountResponse{
			StatusCode:   0,
			StatusMsg:    "success",
			CommentCount: count,
		}, nil
	}
	count, err = query.Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(videoId))).Count()
	if err != nil {
		klog.Infof("video comment count query failed err:%v", err)
		return nil, err
	}
	err = cache.Client.Set(ctx, fmt.Sprintf("interaction.rpc:videoCommentCount:%d", videoId), count, time.Hour).Err()
	if err != nil {
		klog.Infof("cache failed cache interaction.rpc:videoCommentCount")
	}
	return &interaction.GetVideoCommentCountResponse{
		StatusCode:   0,
		StatusMsg:    "success",
		CommentCount: count,
	}, nil
}

// IsFavorite implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) IsFavorite(ctx context.Context, req *interaction.IsFavoriteRequest) (resp *interaction.IsFavoriteResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	videoId := req.VideoId
	favourite, _ := query.Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId)), query.UserFavourite.Status.Eq(1)).First()
	if favourite == nil {
		return &interaction.IsFavoriteResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			IsFavorite: false,
		}, nil
	}
	return &interaction.IsFavoriteResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		IsFavorite: true,
	}, nil
}

// GetCommentById implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetCommentById(ctx context.Context, req *interaction.GetCommentByIdRequest) (resp *interaction.GetCommentByIdResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	commentId := req.CommentId
	comment, err := query.Q.Comment.Where(query.Comment.ID.Eq(uint(commentId))).First()
	if comment == nil {
		klog.Infof("commentId is not existed")
		return nil, err
	}

	return &interaction.GetCommentByIdResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		CommentId:  commentId,
		UserId:     int64(comment.FromUserId),
		Content:    comment.Content,
		VideoId:    int64(comment.ToVideoId),
	}, nil
}

// 校验req是否为空，空值直接返回err
func checkReq(req interface{}) error {
	if req == nil {
		klog.Warnf("req is nil please check other service")
		return errors.New("req is nil please check other service")
	}
	return nil
}

// GetInteractionInfo implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetInteractionInfo(ctx context.Context, req *interaction.GetInteractionInfoRequest) (resp *interaction.GetInteractionInfoResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	//喜欢数
	favoriteCount, err := query.Q.UserFavourite.Where(query.Q.UserFavourite.UserId.Eq(uint(userId)), query.Q.UserFavourite.Status.Eq(1)).Count()
	if err != nil {
		klog.Infof("query favorite count failed err:%v", err)
		return nil, err
	}
	//作品数
	GetUserVideoIdsResponse, err := BasicsRpcClient.GetUserVideoIds(ctx, &basics.UserVideoIdsRequest{UserId: userId})
	if err != nil {
		klog.Infof("BasicsRpcClient failed err:%v", err)
		return nil, err
	}
	var workCount int64
	var totalFavorited int64
	if GetUserVideoIdsResponse != nil {
		videoIdList := GetUserVideoIdsResponse.VideoIdList
		workCount = int64(len(videoIdList))

		if workCount != 0 {
			videoIdListUInt := make([]uint, 0, len(videoIdList))
			for _, videoId := range videoIdList {
				videoIdListUInt = append(videoIdListUInt, uint(videoId))
			}
			totalFavorited, err = query.Q.UserFavourite.Where(query.Q.UserFavourite.VideoId.In(videoIdListUInt...), query.Q.UserFavourite.Status.Eq(1)).Count()
			if err != nil {
				klog.Infof("query total favorited failed err:%v", err)
				return nil, err
			}
		}
	}
	return &interaction.GetInteractionInfoResponse{FavoriteCount: favoriteCount, WorkCount: workCount, TotalFavorited: totalFavorited}, nil
}
