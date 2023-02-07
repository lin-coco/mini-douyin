package main

import (
	"basics.rpc/kitex_gen/douyin/core"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"log"
	"society.rpc/dal/model"
	"society.rpc/dal/query"
	second "society.rpc/kitex_gen/douyin/extra/second"
	"time"
)

// SocietyServiceImpl implements the last service interface defined in the IDL.
type SocietyServiceImpl struct{}

// ConcernAction implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) ConcernAction(ctx context.Context, req *second.ConcernActionRequest) (resp *second.ConcernActionResponse, err error) {
	fromUserId := req.FromUserId
	toUserId := req.ToUserId
	if fromUserId == toUserId {
		log.Printf("concern failed because fromUserId:%d = toUserId:%d", fromUserId, toUserId)
		return nil, errors.New("concern failed because fromUserId = toUserId")
	}
	relation, _ := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).First()
	if relation != nil {
		_, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).Update(query.Relation.RelType, 1)
		if err != nil {
			log.Printf("concern failed fromUserId:%d toUserId:%d err:%v", fromUserId, toUserId, err)
			return nil, err
		}
	} else {
		err := Q.Relation.Create(&model.Relation{FromUserId: uint(fromUserId), ToUserId: uint(toUserId)})
		if err != nil {
			log.Printf("concern failed fromUserId:%d toUserId:%d err:%v", fromUserId, toUserId, err)
			return nil, err
		}
	}

	return &second.ConcernActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// CancelConcernAction implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) CancelConcernAction(ctx context.Context, req *second.CancelConcernActionRequest) (resp *second.CancelConcernActionResponse, err error) {
	fromUserId := req.FromUserId
	toUserId := req.ToUserId
	relation, _ := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).First()
	if relation == nil {
		log.Printf("has no concerned fromUserId:%d toUserId:%d", fromUserId, toUserId)
	} else {
		_, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).Update(query.Relation.RelType, 0)
		if err != nil {
			log.Printf("calcel concern failed fromUserId:%d toUserId:%d err:%v", fromUserId, toUserId, err)
			return nil, err
		}
	}
	return &second.CancelConcernActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// FollowList implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) FollowList(ctx context.Context, req *second.FollowListRequest) (resp *second.FollowListResponse, err error) {
	userId := req.UserId
	relations, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	if err != nil {
		log.Printf("query follow failed userId:%d err:%v", userId, err)
		return nil, err
	}
	toUserIds := make([]int64, 0, len(relations))
	for _, relation := range relations {
		toUserIds = append(toUserIds, int64(relation.ToUserId))
	}
	res, err := BasicsService.GetUserListByIds(ctx, &core.GetUserListByIdsRequest{UserIdList: toUserIds})
	if err != nil {
		log.Printf("BasicsService run failed err:%v", err)
		return nil, err
	}
	userList := res.UserList
	users := make([]*second.User, 0, len(userList))
	for _, user := range userList {
		followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follow count failed err:%v", user.Id, err)
		}
		followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follower count failed err:%v", user.Id, err)
		}
		users = append(users, &second.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      true,
		})
	}
	return &second.FollowListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// FollowerList implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) FollowerList(ctx context.Context, req *second.FollowerListRequest) (resp *second.FollowerListResponse, err error) {
	userId := req.UserId
	relations, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	if err != nil {
		log.Printf("query follower failed userId:%d err:%v", userId, err)
		return nil, err
	}
	if len(relations) == 0 {
		return &second.FollowerListResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*second.User, 0, 0),
		}, nil
	}
	FromUserIds := make([]int64, 0, len(relations))
	for _, relation := range relations {
		FromUserIds = append(FromUserIds, int64(relation.FromUserId))
	}
	res, err := BasicsService.GetUserListByIds(ctx, &core.GetUserListByIdsRequest{UserIdList: FromUserIds})
	if err != nil {
		log.Printf("BasicsService run failed err:%v", err)
		return nil, err
	}
	userList := res.UserList
	users := make([]*second.User, 0, len(userList))

	for _, user := range userList {
		followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follow count failed err:%v", user.Id, err)
		}
		followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follower count failed err:%v", user.Id, err)
		}
		users = append(users, &second.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      false,
		})
	}
	return &second.FollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// FriendList implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) FriendList(ctx context.Context, req *second.FriendListRequest) (resp *second.FriendListResponse, err error) {
	userId := req.UserId
	//先找我关注的
	relations, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	if err != nil {
		log.Printf("query follow failed userId:%d err:%v", userId, err)
		return nil, err
	}
	toUserIds := make([]uint, 0, len(relations))
	for _, relation := range relations {
		toUserIds = append(toUserIds, relation.ToUserId)
	}
	//如果我关注的人也关注我即为好友
	friends, err := Q.Relation.Where(query.Relation.FromUserId.In(toUserIds...), query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	if err != nil {
		log.Printf("query friend failed userId:%d err:%v", userId, err)
	}
	friendIds := make([]int64, 0, len(friends))
	for _, friend := range friends {
		friendIds = append(friendIds, int64(friend.FromUserId))
	}
	if len(friendIds) == 0 {
		return &second.FriendListResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*second.User, 0, 0),
		}, nil
	}
	res, err := BasicsService.GetUserListByIds(ctx, &core.GetUserListByIdsRequest{UserIdList: friendIds})
	if err != nil {
		log.Printf("BasicsService run failed err:%v", err)
		return nil, err
	}
	friendList := res.UserList
	users := make([]*second.User, 0, len(friendList))
	for _, friend := range friendList {
		followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(friend.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follow count failed err:%v", friend.Id, err)
		}
		followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(friend.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follower count failed err:%v", friend.Id, err)
		}
		users = append(users, &second.User{
			Id:            friend.Id,
			Name:          friend.Name,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      true,
		})
	}

	return &second.FriendListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// SocietyInfo implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) SocietyInfo(ctx context.Context, req *second.SocietyInfoRequest) (resp *second.SocietyInfoResponse, err error) {
	myId := req.MyId
	userId := req.UserId
	followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	if err != nil {
		log.Printf("userId:%d query follow count failed err:%v", userId, err)
		return nil, err
	}
	followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	if err != nil {
		log.Printf("userId:%d query follower count failed err:%v", userId, err)
		return nil, err
	}
	count, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(myId)), query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	if err != nil {
		log.Printf("myId:%d userId:%d query relation count failed err:%v", myId, userId, err)
		return nil, err
	}
	isFollow := false
	if count > 0 || myId == userId {
		isFollow = true
	}
	return &second.SocietyInfoResponse{
		StatusCode:    0,
		StatusMsg:     "success",
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}, nil
}

// MessageChat implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) MessageChat(ctx context.Context, req *second.MessageChatRequest) (resp *second.MessageChatResponse, err error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	myUserId := req.MyUserId
	friendUserId := req.FriendUserId
	if myUserId == friendUserId {
		return nil, errors.New("myUserId = friendUserId error")
	}
	result, err := RedisDB.LRange(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			klog.Infof("redis not cache society.rpc:messagechat")
		} else {
			klog.Infof("redis failed society.rpc:messagechat")
		}
	} else {
		messages := make([]*second.Message, 0, len(result))
		for _, ms := range result {
			m := new(second.Message)
			err := json.Unmarshal([]byte(ms), m)
			if err != nil {
				klog.Infof("message unmarshal failed")
				continue
			}
			messages = append(messages, m)
		}
		return &second.MessageChatResponse{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: messages,
		}, nil
	}
	result, err = RedisDB.LRange(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", friendUserId, myUserId), 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			klog.Infof("redis not cache society.rpc:messagechat")
		} else {
			klog.Infof("redis failed society.rpc:messagechat")
		}
	} else {
		messages := make([]*second.Message, 0, len(result))
		for _, ms := range result {
			m := new(second.Message)
			err := json.Unmarshal([]byte(ms), m)
			if err != nil {
				klog.Infof("message unmarshal failed")
				continue
			}
			messages = append(messages, m)
		}
		return &second.MessageChatResponse{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: messages,
		}, nil
	}

	messageChats, err := Q.MessageChat.Where(query.MessageChat.FromUserId.In(uint(myUserId), uint(friendUserId)), query.MessageChat.ToUserId.In(uint(myUserId), uint(friendUserId))).Order(query.MessageChat.CreatedAt).Find()
	if err != nil {
		log.Printf("myId:%d friendId:%d query message chat failed err:%v", myUserId, friendUserId, err)
		return nil, err
	}

	if len(messageChats) == 0 {
		//RedisDB.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:fromUserId%dToUserId%d", myUserId, friendUserId), time.Hour)
		return &second.MessageChatResponse{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: make([]*second.Message, 0, 0),
		}, nil
	}
	messages := make([]*second.Message, 0, len(messageChats))
	redMsgs := make([]string, 0, len(messageChats))
	for _, chat := range messageChats {
		createTime := chat.CreatedAt.Format("2006-01-02 15:04")
		msg := &second.Message{
			Id:         int64(chat.ID),
			FromUserId: int64(chat.FromUserId),
			ToUserId:   int64(chat.ToUserId),
			Content:    chat.MsgContent,
			CreateTime: &createTime,
		}
		messages = append(messages, msg)
		msgByte, _ := json.Marshal(msg)
		redMsgs = append(redMsgs, string(msgByte))
	}
	RedisDB.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), redMsgs, time.Hour)

	return &second.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		MessageList: messages,
	}, nil
}

// MessageSend implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) MessageSend(ctx context.Context, req *second.MessageSendRequest) (resp *second.MessageSendResponse, err error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	myUserId := req.MyUserId
	friendUserId := req.FriendUserId
	content := req.Content
	if myUserId == friendUserId {
		return nil, errors.New("myUserId = friendUserId error")
	}
	m := &model.MessageChat{
		MsgContent: content,
		FromUserId: uint(myUserId),
		ToUserId:   uint(friendUserId)}
	err = Q.MessageChat.Create(m)
	if err != nil {
		log.Printf("message send failed FromUserId:%d ToUserId:%d MsgContent:%s", myUserId, friendUserId, content)
		return nil, err
	}
	//redis
	result, _ := RedisDB.Exists(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId)).Result()
	if result > 0 {
		createTime := m.CreatedAt.Format("2006-01-02 15:04")
		//存在
		sM := &second.Message{
			Id:         int64(m.ID),
			FromUserId: int64(m.FromUserId),
			ToUserId:   int64(m.ToUserId),
			Content:    m.MsgContent,
			CreateTime: &createTime,
		}
		bytes, _ := json.Marshal(sM)
		RedisDB.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), string(bytes))
	}
	result, _ = RedisDB.Exists(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", friendUserId, myUserId)).Result()
	if result > 0 {
		createTime := m.CreatedAt.Format("2006-01-02 15:04")
		//存在
		sM := &second.Message{
			Id:         int64(m.ID),
			FromUserId: int64(m.FromUserId),
			ToUserId:   int64(m.ToUserId),
			Content:    m.MsgContent,
			CreateTime: &createTime,
		}
		bytes, _ := json.Marshal(sM)
		RedisDB.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), string(bytes))
	}
	return &second.MessageSendResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}
