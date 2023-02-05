package main

import (
	"basics.rpc/dal/model"
	"basics.rpc/dal/query"
	core "basics.rpc/kitex_gen/douyin/core"
	"basics.rpc/utils"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"log"
	"time"
)

// BasicsServiceImpl implements the last service interface defined in the IDL.
type BasicsServiceImpl struct{}

// GetUserInfoById implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetUserInfoById(ctx context.Context, req *core.GetUserRequest) (resp *core.GetUserResponse, err error) {
	userId := req.UserId
	user, err := query.Q.User.Where(query.User.ID.Eq(uint(userId))).First()
	if err != nil {
		log.Printf("query failed userId:%d, err:%v", userId, err)
		return nil, err
	}
	return &core.GetUserResponse{
		Id:   int64(user.ID),
		Name: user.Name,
	}, nil
}

// CreateUser implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) CreateUser(ctx context.Context, req *core.CreateUserRequest) (resp *core.CreateUserResponse, err error) {
	// TODO: Your code here...
	username := req.Username
	password := req.Password
	first, _ := query.User.Where(query.User.Name.Eq(username)).First()
	if first != nil {
		log.Printf("username has existed %s", username)
		return nil, errors.New("username has existed " + username)
	}
	//密码加密
	passwordHash, err := utils.PasswordHash(PwdKey + password)
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
	return &core.CreateUserResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Id:         int64(user.ID),
	}, nil
}

// GetVideo implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideo(ctx context.Context, req *core.GetVideoRequest) (resp *core.GetVideoResponse, err error) {
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
	videos, err := query.Q.Video.Order(query.Video.CreatedAt).Where(query.Video.CreatedAt.Lt(latest)).Limit(30).Find()
	if err != nil {
		log.Printf("getvideo failed time:%v err:%v", latest, err)
		return nil, err
	}

	videoList := make([]*core.Video, 0, len(videos))
	nextTime := latest
	for i := 0; i < len(videos); i++ {
		user, err := query.Q.User.Where(query.User.ID.Eq(videos[i].UserId)).First()
		if err != nil {
			log.Printf("query User by Id failed err:%v", err)
			return nil, err
		}
		videoList = append(videoList, &core.Video{
			Id:       int64(videos[i].ID),
			User:     &core.User{Name: user.Name, Id: int64(user.ID)},
			PlayUrl:  videos[i].PlayUrl,
			CoverUrl: videos[i].CoverUrl,
			Title:    videos[i].Title,
		})
		if videos[i].CreatedAt.Before(nextTime) {
			nextTime = videos[i].CreatedAt
		}
	}

	return &core.GetVideoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
		NextTime:   nextTime.Unix(),
	}, nil
}

// UploadVideo implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) UploadVideo(ctx context.Context, req *core.UploadVideoRequest) (resp *core.UploadVideoResponse, err error) {
	// TODO: Your code here...
	data := req.Data
	title := req.Title
	userId := req.UserId

	//cmdIn, cmdInWriter := io.Pipe()
	//_, err = io.Copy(cmdInWriter, bytes.NewReader(data))
	//if err != nil {
	//	log.Printf("coverurl failed err:%v", err)
	//	return nil, err
	//}
	//err = cmdInWriter.Close()
	//if err != nil {
	//	log.Printf("coverurl failed err:%v", err)
	//	return nil, err
	//}
	////封面图片字节数据
	//cmd := exec.Command("ffmpeg", "-i", "-vframes", "1", "-q:v", "2", "-f", "image2", "pipe:1")
	//cmd.Stdin = cmdIn
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err = cmd.Run()
	//if err != nil {
	//	log.Printf("coverurl failed err:%v", err)
	//	return nil, err
	//}
	//coverData := out.Bytes()

	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	time.FixedZone("Bei jing Time", secondsEastOfUTC)
	now := time.Now()
	fileName := fmt.Sprintf("uservideo/%d%s%d-%d-%d|%d:%d:%d|%d.mp4", userId, title, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.UnixMicro())
	//coverFileName := fmt.Sprintf("uservideo/%d%s%d.%d.%d.jpg", userId, title, now.Year(), now.Month(), now.Day())
	//fileName := "uservideo/" + strconv.Itoa(userId) + title + string(now.Year()) + "." + string(now.Month()) + "." + string(now.Day()) + ".mp4"
	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//b := atomic.Bool{}
	//func() {
	err = OSSBucket.PutObject(fileName, bytes.NewReader(data))
	if err != nil {
		log.Printf("upload failed err:%v", err)
		//b.Store(true)
		return nil, err
	}
	//wg.Done()
	//}()
	//go func() {
	//err = OSSBucket.PutObject(coverFileName, bytes.NewReader(coverData))
	//if err != nil {
	//	log.Printf("upload failed err:%v", err)
	//	//b.Store(true)
	//	return nil, err
	//}
	//wg.Done()
	//}()
	//wg.Wait()
	//if b.Load() {
	//	log.Printf("upload failed err:%v", err)
	//	return nil, err
	//}

	err = query.Q.Video.Create(&model.Video{
		UserId:  uint(userId),
		PlayUrl: OSSBaseUrl + fileName,
		//CoverUrl: OSSBaseUrl + coverFileName,
		Title: title,
	})
	if err != nil {
		log.Printf("create video failed err:%v", err)
		return nil, err
	}
	return &core.UploadVideoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// GetVideosByUserId implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideosByUserId(ctx context.Context, req *core.GetVideosByUserIdRequest) (resp *core.GetVideosByUserIdResponse, err error) {
	userId := req.UserId
	videos, err := query.Q.Video.Where(query.Video.UserId.Eq(uint(userId))).Find()
	if err != nil {
		log.Printf("query failed err:%v", err)
		return nil, err
	}
	videoList := make([]*core.Video, 0, len(videos))
	for _, video := range videos {
		user, err := query.Q.User.Where(query.User.ID.Eq(uint(userId))).First()
		if err != nil || user == nil {
			log.Printf("query user failed userId:%d err:%v", userId, err)
			return nil, err
		}
		videoList = append(videoList, &core.Video{
			Id:       int64(video.ID),
			User:     &core.User{Id: int64(user.ID), Name: user.Name},
			PlayUrl:  video.PlayUrl,
			CoverUrl: video.CoverUrl,
			Title:    video.Title,
		})
	}

	return &core.GetVideosByUserIdResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}, nil
}

// CheckUser implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) CheckUser(ctx context.Context, req *core.CheckUserRequest) (resp *core.CheckUserResponse, err error) {
	username := req.Username
	password := req.Password

	user, _ := query.Q.User.Where(query.User.Name.Eq(username)).First()
	if user == nil {
		log.Printf("query user failed username:%s err:%v", username, err)
		return &core.CheckUserResponse{
			StatusCode: 1,
			StatusMsg:  "username is not effective",
		}, errors.New("username is not effective")
	}
	if !utils.PasswordVerify(PwdKey+password, user.Password) {
		return &core.CheckUserResponse{
			StatusCode: 1,
			StatusMsg:  "password is not effective",
		}, errors.New("username is not effective")
	}

	return &core.CheckUserResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(user.ID),
	}, nil
}

// GetVideoListByIds implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetVideoListByIds(ctx context.Context, req *core.GetVideoListByIdsRequest) (resp *core.GetVideoListByIdsResponse, err error) {
	if req == nil {
		klog.Infof("req is nil failed")
		return nil, errors.New("req is nil failed")
	}
	videoIdList := req.VideoIdList
	videoList := make([]*core.Video, 0, len(videoIdList))
	for _, videoId := range videoIdList {
		video, err := query.Q.Video.Where(query.Video.ID.Eq(uint(videoId))).First()
		if err != nil || video == nil {
			log.Printf("one no discover videoId:%d err:%v", videoId, err)
			continue
		}
		user, err := query.Q.User.Where(query.User.ID.Eq(video.UserId)).First()
		if err != nil || user == nil {
			log.Printf("one no discover userId:%d err:%v", video.UserId, err)
			videoList = append(videoList, &core.Video{Id: videoId, PlayUrl: video.PlayUrl, CoverUrl: video.CoverUrl, Title: video.Title})
			continue
		}
		videoList = append(videoList, &core.Video{
			Id:       videoId,
			User:     &core.User{Id: int64(user.ID), Name: user.Name},
			PlayUrl:  video.PlayUrl,
			CoverUrl: video.CoverUrl,
			Title:    video.Title})
	}
	return &core.GetVideoListByIdsResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}, nil
}

// GetUserListByIds implements the BasicsServiceImpl interface.
func (s *BasicsServiceImpl) GetUserListByIds(ctx context.Context, req *core.GetUserListByIdsRequest) (resp *core.GetUserListByIdsResponse, err error) {
	if req == nil {
		return &core.GetUserListByIdsResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*core.User, 0, 0),
		}, nil
	}
	userIdList := req.UserIdList
	if userIdList == nil || len(userIdList) == 0 {
		return &core.GetUserListByIdsResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*core.User, 0, 0),
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
	userList := make([]*core.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, &core.User{Id: int64(user.ID), Name: user.Name})
	}
	return &core.GetUserListByIdsResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
	}, nil
}
