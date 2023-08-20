package basics_rpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"log"
	"mini-douyin/internal/pkg/dal/model"
	"mini-douyin/internal/pkg/dal/query"
	"mini-douyin/internal/pkg/kitex_gen/douyin/basics"
	"mini-douyin/internal/pkg/utils"
	"mini-douyin/pkg/cache"
	"mini-douyin/pkg/oss"
	"os/exec"
	"time"
)

// BasicsServiceImpl implements the last service interface defined in the IDL.
type BasicsServiceImpl struct{}

// GetUserInfoById implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetUserInfoById(ctx context.Context, req *basics.GetUserRequest) (resp *basics.GetUserResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	user, err := query.Q.User.Where(query.User.ID.Eq(uint(userId))).First()
	if err != nil {
		log.Printf("query failed userId:%d, err:%v", userId, err)
		return nil, err
	}
	return &basics.GetUserResponse{
		Id:   int64(user.ID),
		Name: user.Name,
	}, nil
}

// CreateUser implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) CreateUser(ctx context.Context, req *basics.CreateUserRequest) (resp *basics.CreateUserResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	username := req.Username
	password := req.Password
	val := cache.Client.LRange(ctx, "basics.rpc:userexist", 0, -1).Val()
	for _, s := range val {
		if s == username {
			return nil, errors.New("username has existed " + username)
		}
	}
	first, _ := query.User.Where(query.User.Name.Eq(username)).First()
	if first != nil {
		log.Printf("username has existed %s", username)
		//cache 缓存优化
		err := cache.Client.LPush(ctx, fmt.Sprintf("basics.rpc:userexist"), username).Err()
		cache.Client.Expire(ctx, fmt.Sprintf("basics.rpc:userexist"), time.Hour)
		if err != nil {
			klog.Infof("cache cache basics.rpc:userexist failed")
		}
		return nil, errors.New("username has existed " + username)
	}
	//密码加密
	passwordHash, err := utils.PasswordHash(Config.PwdKey + password)
	if err != nil {
		log.Printf("hash password failed err:%v", err)
		return nil, err
	}
	user := &model.User{Name: username, Password: passwordHash}
	err = query.Q.User.Create(user)
	if err != nil {
		log.Fatalf("create failed username:%d, err:%v", username, err)
		return nil, err
	}
	err = cache.Client.LRem(ctx, "basics.rpc:usernotexist", 1, username).Err()
	if err != nil {
		klog.Infof("cache rem usernotexist failed")
	}
	return &basics.CreateUserResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Id:         int64(user.ID),
	}, nil
}

// GetVideo implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideo(ctx context.Context, req *basics.GetVideoRequest) (resp *basics.GetVideoResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	var latestTime int64
	if req != nil {
		latestTime = req.LatestTime //秒级时间戳
	}
	//东八区时间对象
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	time.FixedZone("Bei jing Time", secondsEastOfUTC)
	var latest time.Time
	if latestTime == 0 {
		latest = time.Now()
	} else {
		latest = time.Unix(latestTime, 0)
	}
	videos, err := query.Q.Video.Where(query.Video.CreatedAt.Lt(latest)).Order(query.Video.CreatedAt.Desc()).Limit(30).Find()
	if err != nil {
		log.Printf("getvideo failed time:%v err:%v", latest, err)
		return nil, err
	}

	videoList := make([]*basics.Video, 0, len(videos))
	var nextTime int64
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].CreatedAt.UnixMilli()
	}
	for i := 0; i < len(videos); i++ {
		user, err := query.Q.User.Where(query.User.ID.Eq(videos[i].UserId)).First()
		if err != nil {
			log.Printf("query User by Id failed err:%v", err)
			return nil, err
		}
		videoList = append(videoList, &basics.Video{
			Id:       int64(videos[i].ID),
			User:     &basics.User{Name: user.Name, Id: int64(user.ID)},
			PlayUrl:  videos[i].PlayUrl,
			CoverUrl: videos[i].CoverUrl,
			Title:    videos[i].Title,
		})
		//if videos[i].CreatedAt.Before(nextTime) {
		//	nextTime = videos[i].CreatedAt
		//}
	}

	return &basics.GetVideoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
		NextTime:   nextTime,
	}, nil
}

// UploadVideo implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) UploadVideo(ctx context.Context, req *basics.UploadVideoRequest) (resp *basics.UploadVideoResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	data := req.Data
	title := req.Title
	userId := req.UserId

	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	time.FixedZone("Bei jing Time", secondsEastOfUTC)
	now := time.Now()
	//上传视频
	fileName := fmt.Sprintf("uservideo/%d%s%d-%d-%d|%d:%d:%d|%d.mp4", userId, title, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.UnixMicro())
	err = oss.OSSBucket.PutObject(fileName, bytes.NewReader(data))
	if err != nil {
		log.Printf("upload failed err:%v", err)
		return nil, err
	}
	//上传封面
	cmd := exec.Command("ffmpeg", "-i", oss.Config.OSS.BaseUrl+fileName, "-vframes", "1", "-q:v", "2", "-f", "image2", "pipe:1")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Decode the image data from the FFmpeg output
	imgData := out.Bytes()
	if err != nil {
		fmt.Println(err)
		return
	}
	coverFileName := fmt.Sprintf("uservideo/%d%s%d-%d-%d|%d:%d:%d|%d.jpg", userId, title, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.UnixMicro())
	err = oss.OSSBucket.PutObject(coverFileName, bytes.NewReader(imgData))
	if err != nil {
		log.Printf("upload failed err:%v", err)
		return nil, err
	}
	err = query.Q.Video.Create(&model.Video{
		UserId:   uint(userId),
		PlayUrl:  oss.Config.OSS.BaseUrl + fileName,
		CoverUrl: oss.Config.OSS.BaseUrl + coverFileName,
		Title:    title,
	})
	if err != nil {
		log.Printf("create video failed err:%v", err)
		return nil, err
	}
	return &basics.UploadVideoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// GetVideosByUserId implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideosByUserId(ctx context.Context, req *basics.GetVideosByUserIdRequest) (resp *basics.GetVideosByUserIdResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	count, err := query.Q.User.Where(query.User.ID.Eq(uint(userId))).Count()
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("userInfo is not found")
	}
	videos, err := query.Q.Video.Where(query.Video.UserId.Eq(uint(userId))).Find()
	if err != nil {
		log.Printf("query failed err:%v", err)
		return nil, err
	}
	videoList := make([]*basics.Video, 0, len(videos))
	for _, video := range videos {
		user, err := query.Q.User.Where(query.User.ID.Eq(uint(userId))).First()
		if err != nil || user == nil {
			log.Printf("query user failed userId:%d err:%v", userId, err)
			return nil, err
		}
		videoList = append(videoList, &basics.Video{
			Id:       int64(video.ID),
			User:     &basics.User{Id: int64(user.ID), Name: user.Name},
			PlayUrl:  video.PlayUrl,
			CoverUrl: video.CoverUrl,
			Title:    video.Title,
		})
	}

	return &basics.GetVideosByUserIdResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}, nil
}

// CheckUser implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) CheckUser(ctx context.Context, req *basics.CheckUserRequest) (resp *basics.CheckUserResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	username := req.Username
	password := req.Password

	val := cache.Client.LRange(ctx, "basics.rpc:usernotexist", 0, -1).Val()
	for _, s := range val {
		if s == username {
			return &basics.CheckUserResponse{
				StatusCode: 1,
				StatusMsg:  "username is not effective",
			}, errors.New("username is not effective")
		}
	}
	user, _ := query.Q.User.Where(query.User.Name.Eq(username)).First()
	if user == nil {
		klog.Infof("query user failed username:%s err:%v", username, err)
		err := cache.Client.LPush(ctx, fmt.Sprintf("basics.rpc:usernotexist"), username).Err()
		cache.Client.Expire(ctx, fmt.Sprintf("basics.rpc:usernotexist"), time.Hour)
		if err != nil {
			klog.Infof("cache cache basics.rpc:usernotexist failed")
		}
		return &basics.CheckUserResponse{
			StatusCode: 1,
			StatusMsg:  "username is not effective",
		}, errors.New("username is not effective")
	}
	if !utils.PasswordVerify(Config.PwdKey+password, user.Password) {
		return &basics.CheckUserResponse{
			StatusCode: 1,
			StatusMsg:  "password is not effective",
		}, errors.New("password is not effective")
	}

	return &basics.CheckUserResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(user.ID),
	}, nil
}

// GetVideoListByIds implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideoListByIds(ctx context.Context, req *basics.GetVideoListByIdsRequest) (resp *basics.GetVideoListByIdsResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	videoIdList := req.VideoIdList
	videoList := make([]*basics.Video, 0, len(videoIdList))
	for _, videoId := range videoIdList {
		video, err := query.Q.Video.Where(query.Video.ID.Eq(uint(videoId))).First()
		if err != nil || video == nil {
			log.Printf("one no discover videoId:%d err:%v", videoId, err)
			continue
		}
		user, err := query.Q.User.Where(query.User.ID.Eq(video.UserId)).First()
		if err != nil || user == nil {
			log.Printf("one no discover userId:%d err:%v", video.UserId, err)
			videoList = append(videoList, &basics.Video{Id: videoId, PlayUrl: video.PlayUrl, CoverUrl: video.CoverUrl, Title: video.Title})
			continue
		}
		videoList = append(videoList, &basics.Video{
			Id:       videoId,
			User:     &basics.User{Id: int64(user.ID), Name: user.Name},
			PlayUrl:  video.PlayUrl,
			CoverUrl: video.CoverUrl,
			Title:    video.Title})
	}
	return &basics.GetVideoListByIdsResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}, nil
}

// GetUserListByIds implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetUserListByIds(ctx context.Context, req *basics.GetUserListByIdsRequest) (resp *basics.GetUserListByIdsResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return &basics.GetUserListByIdsResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*basics.User, 0, 0),
		}, nil
	}
	userIdList := req.UserIdList
	if userIdList == nil || len(userIdList) == 0 {
		return &basics.GetUserListByIdsResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*basics.User, 0, 0),
		}, nil
	}
	userIdList2 := make([]uint, 0, len(userIdList))
	for _, userId := range userIdList {
		userIdList2 = append(userIdList2, uint(userId))
	}
	users, err := query.Q.User.Where(query.User.ID.In(userIdList2...)).Find()
	if err != nil {
		log.Printf("query users failed err:%v", err)
		return nil, err
	}
	userList := make([]*basics.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, &basics.User{Id: int64(user.ID), Name: user.Name})
	}
	return &basics.GetUserListByIdsResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
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

// GetVideoInfoById implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideoInfoById(ctx context.Context, req *basics.GetVideoByIdRequest) (resp *basics.GetVideoByIdResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	videoId := req.VideoId
	video, err := query.Q.Video.Where(query.Video.ID.Eq(uint(videoId))).First()
	if err != nil {
		klog.Infof("query failed videoId:%d, err:%v", videoId, err)
		return nil, err
	}
	return &basics.GetVideoByIdResponse{
		Id:    int64(video.ID),
		Title: video.Title,
	}, nil
}

// GetVideoCount implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideoCount(ctx context.Context, req *basics.VideoCountRequest) (resp *basics.VideoCountResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	count, err := query.Video.Where(query.Video.UserId.Eq(uint(userId))).Count()
	if err != nil {
		klog.Infof("query failed user video count:%d, err:%v", userId, err)
		return nil, err
	}
	return &basics.VideoCountResponse{WorkCount: count}, nil
}

// GetUserVideoIds implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetUserVideoIds(ctx context.Context, req *basics.UserVideoIdsRequest) (resp *basics.UserVideoIdsResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	videos, err := query.Video.Select(query.Video.ID).Where(query.Video.UserId.Eq(uint(userId))).Find()
	if err != nil {
		klog.Infof("query failed user video ids:%d, err:%v", userId, err)
		return nil, err
	}
	videoIds := make([]int64, 0, len(videos))
	if len(videos) > 0 {
		for _, video := range videos {
			videoIds = append(videoIds, int64(video.ID))
		}
	}
	return &basics.UserVideoIdsResponse{VideoIdList: videoIds}, nil
}
