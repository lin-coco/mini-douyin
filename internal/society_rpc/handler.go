package society_rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"log"
	"mini-douyin/internal/pkg/dal/model"
	"mini-douyin/internal/pkg/dal/query"
	"mini-douyin/internal/pkg/kitex_gen/douyin/basics"
	"mini-douyin/internal/pkg/kitex_gen/douyin/society"
	"mini-douyin/pkg/cache"
	"time"
)

// SocietyServiceImpl implements the last service interface defined in the IDL.
type SocietyServiceImpl struct{}

// ConcernAction implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) ConcernAction(ctx context.Context, req *society.ConcernActionRequest) (resp *society.ConcernActionResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	fromUserId := req.FromUserId
	toUserId := req.ToUserId
	if fromUserId == toUserId {
		log.Printf("concern failed because fromUserId:%d = toUserId:%d", fromUserId, toUserId)
		return nil, errors.New("concern failed because fromUserId = toUserId")
	}
	relation, _ := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).First()
	if relation != nil {
		_, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).Update(query.Relation.RelType, 1)
		if err != nil {
			log.Printf("concern failed fromUserId:%d toUserId:%d err:%v", fromUserId, toUserId, err)
			return nil, err
		}
	} else {
		err := query.Q.Relation.Create(&model.Relation{FromUserId: uint(fromUserId), ToUserId: uint(toUserId)})
		if err != nil {
			log.Printf("concern failed fromUserId:%d toUserId:%d err:%v", fromUserId, toUserId, err)
			return nil, err
		}
	}

	return &society.ConcernActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// CancelConcernAction implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) CancelConcernAction(ctx context.Context, req *society.CancelConcernActionRequest) (resp *society.CancelConcernActionResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	fromUserId := req.FromUserId
	toUserId := req.ToUserId
	relation, _ := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).First()
	if relation == nil {
		log.Printf("has no concerned fromUserId:%d toUserId:%d", fromUserId, toUserId)
	} else {
		_, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).Update(query.Relation.RelType, 0)
		if err != nil {
			log.Printf("calcel concern failed fromUserId:%d toUserId:%d err:%v", fromUserId, toUserId, err)
			return nil, err
		}
	}
	return &society.CancelConcernActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// FollowList implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) FollowList(ctx context.Context, req *society.FollowListRequest) (resp *society.FollowListResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	relations, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	if err != nil {
		log.Printf("query follow failed userId:%d err:%v", userId, err)
		return nil, err
	}
	toUserIds := make([]int64, 0, len(relations))
	for _, relation := range relations {
		toUserIds = append(toUserIds, int64(relation.ToUserId))
	}
	res, err := BasicsRpcClient.GetUserListByIds(ctx, &basics.GetUserListByIdsRequest{UserIdList: toUserIds})
	if err != nil {
		log.Printf("BasicsRpcClient run failed err:%v", err)
		return nil, err
	}
	userList := res.UserList
	users := make([]*society.User, 0, len(userList))
	for _, user := range userList {
		followCount, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follow count failed err:%v", user.Id, err)
		}
		followerCount, err := query.Q.Relation.Where(query.Relation.ToUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follower count failed err:%v", user.Id, err)
		}
		users = append(users, &society.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      true,
		})
	}
	return &society.FollowListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// FollowerList implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) FollowerList(ctx context.Context, req *society.FollowerListRequest) (resp *society.FollowerListResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	var relations []*model.Relation
	if req.PageNo == 0 && req.PageSize == 0 {
		relations, err = query.Q.Relation.Where(query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Order(query.Relation.CreatedAt).Find()
		if err != nil {
			log.Printf("query follower failed userId:%d err:%v", userId, err)
			return nil, err
		}
	} else {
		relations, _, err = query.Q.Relation.Where(query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Order(query.Relation.CreatedAt).FindByPage(int(req.PageNo), int(req.PageSize))
		if err != nil {
			log.Printf("query follower failed userId:%d err:%v", userId, err)
			return nil, err
		}
	}

	if len(relations) == 0 {
		return &society.FollowerListResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*society.User, 0, 0),
		}, nil
	}
	FromUserIds := make([]int64, 0, len(relations))
	for _, relation := range relations {
		FromUserIds = append(FromUserIds, int64(relation.FromUserId))
	}
	res, err := BasicsRpcClient.GetUserListByIds(ctx, &basics.GetUserListByIdsRequest{UserIdList: FromUserIds})
	if err != nil {
		log.Printf("BasicsRpcClient run failed err:%v", err)
		return nil, err
	}
	userList := res.UserList
	users := make([]*society.User, 0, len(userList))

	for _, user := range userList {
		followCount, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follow count failed err:%v", user.Id, err)
		}
		followerCount, err := query.Q.Relation.Where(query.Relation.ToUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follower count failed err:%v", user.Id, err)
		}
		users = append(users, &society.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      false,
		})
	}
	return &society.FollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// FriendList implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) FriendList(ctx context.Context, req *society.FriendListRequest) (resp *society.FriendListResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	userId := req.UserId
	//先找我关注的
	relations, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	if err != nil {
		log.Printf("query follow failed userId:%d err:%v", userId, err)
		return nil, err
	}
	toUserIds := make([]uint, 0, len(relations))
	for _, relation := range relations {
		toUserIds = append(toUserIds, relation.ToUserId)
	}
	//如果我关注的人也关注我即为好友
	friends, err := query.Q.Relation.Where(query.Relation.FromUserId.In(toUserIds...), query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	if err != nil {
		log.Printf("query friend failed userId:%d err:%v", userId, err)
	}
	friendIds := make([]int64, 0, len(friends))
	for _, friend := range friends {
		friendIds = append(friendIds, int64(friend.FromUserId))
	}
	if len(friendIds) == 0 {
		return &society.FriendListResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			UserList:   make([]*society.User, 0, 0),
		}, nil
	}
	res, err := BasicsRpcClient.GetUserListByIds(ctx, &basics.GetUserListByIdsRequest{UserIdList: friendIds})
	if err != nil {
		log.Printf("BasicsRpcClient run failed err:%v", err)
		return nil, err
	}
	friendList := res.UserList
	users := make([]*society.User, 0, len(friendList))
	for _, friend := range friendList {
		followCount, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(friend.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follow count failed err:%v", friend.Id, err)
		}
		followerCount, err := query.Q.Relation.Where(query.Relation.ToUserId.Eq(uint(friend.Id)), query.Relation.RelType.Eq(1)).Count()
		if err != nil {
			log.Printf("userId:%d query follower count failed err:%v", friend.Id, err)
		}
		users = append(users, &society.User{
			Id:            friend.Id,
			Name:          friend.Name,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      true,
		})
	}

	return &society.FriendListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// SocietyInfo implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) SocietyInfo(ctx context.Context, req *society.SocietyInfoRequest) (resp *society.SocietyInfoResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	myId := req.MyId
	userId := req.UserId
	followCount, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	if err != nil {
		log.Printf("userId:%d query follow count failed err:%v", userId, err)
		return nil, err
	}
	followerCount, err := query.Q.Relation.Where(query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	if err != nil {
		log.Printf("userId:%d query follower count failed err:%v", userId, err)
		return nil, err
	}
	query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(myId)))
	count, err := query.Q.Relation.Where(query.Relation.FromUserId.Eq(uint(myId)), query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	if err != nil {
		log.Printf("myId:%d userId:%d query relation count failed err:%v", myId, userId, err)
		return nil, err
	}
	isFollow := false
	if count > 0 || myId == userId {
		isFollow = true
	}
	return &society.SocietyInfoResponse{
		StatusCode:    0,
		StatusMsg:     "success",
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}, nil
}

// MessageChat implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) MessageChat(ctx context.Context, req *society.MessageChatRequest) (resp *society.MessageChatResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	myUserId := req.MyUserId
	friendUserId := req.FriendUserId
	if myUserId == friendUserId {
		return nil, errors.New("myUserId = friendUserId error")
	}
	result, err := cache.Client.LRange(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			klog.Infof("cache not cache society.rpc:messagechat")
		} else {
			klog.Infof("cache failed society.rpc:messagechat")
		}
	}
	if len(result) != 0 {
		messages := make([]*society.Message, 0, len(result))
		for _, ms := range result {
			m := new(society.Message)
			err := json.Unmarshal([]byte(ms), m)
			if err != nil {
				klog.Infof("message unmarshal failed")
				continue
			}
			messages = append(messages, m)
		}
		return &society.MessageChatResponse{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: messages,
		}, nil
	}
	result, err = cache.Client.LRange(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", friendUserId, myUserId), 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			klog.Infof("cache not cache society.rpc:messagechat")
		} else {
			klog.Infof("cache failed society.rpc:messagechat")
		}
	}
	if len(result) != 0 {
		messages := make([]*society.Message, 0, len(result))
		for _, ms := range result {
			m := new(society.Message)
			err := json.Unmarshal([]byte(ms), m)
			if err != nil {
				klog.Infof("message unmarshal failed")
				continue
			}
			messages = append(messages, m)
		}
		return &society.MessageChatResponse{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: messages,
		}, nil
	}
	var messageChats []*model.MessageChat
	if req.StartTime == 0 && req.EndTime == 0 {
		messageChats, err = query.Q.MessageChat.Where(query.MessageChat.FromUserId.In(uint(myUserId), uint(friendUserId)), query.MessageChat.ToUserId.In(uint(myUserId), uint(friendUserId))).Order(query.MessageChat.CreatedAt).Find()
	} else {
		//startTime 与 endTime 不为零值
		messageChats, err = query.Q.MessageChat.Where(query.MessageChat.FromUserId.In(uint(myUserId), uint(friendUserId)), query.MessageChat.ToUserId.In(uint(myUserId), uint(friendUserId)), query.MessageChat.CreatedAt.Gt(time.Unix(req.StartTime, 0)), query.MessageChat.CreatedAt.Lt(time.Unix(req.EndTime, 0))).Order(query.MessageChat.CreatedAt).Find()
	}
	if err != nil {
		log.Printf("myId:%d friendId:%d query message chat failed err:%v", myUserId, friendUserId, err)
		return nil, err
	}

	if len(messageChats) == 0 {
		//cache.Client.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:fromUserId%dToUserId%d", myUserId, friendUserId), time.Hour)
		return &society.MessageChatResponse{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: make([]*society.Message, 0, 0),
		}, nil
	}
	messages := make([]*society.Message, 0, len(messageChats))
	redMsgs := make([]string, 0, len(messageChats))
	for _, chat := range messageChats {
		createTimeFormat := chat.CreatedAt.Format("2006-01-02 15:04")
		msg := &society.Message{
			Id:               int64(chat.ID),
			FromUserId:       int64(chat.FromUserId),
			ToUserId:         int64(chat.ToUserId),
			Content:          chat.MsgContent,
			CreateTime:       chat.CreatedAt.UnixMilli(),
			CreateTimeFormat: &createTimeFormat,
		}
		messages = append(messages, msg)
		msgByte, _ := json.Marshal(msg)
		redMsgs = append(redMsgs, string(msgByte))
	}
	if len(redMsgs) == 0 {
		klog.Infof("messagechat cache cache no send because list len is 0")
	} else {
		cache.Client.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), redMsgs)
		cache.Client.Expire(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), time.Hour)
	}
	return &society.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		MessageList: messages,
	}, nil
}

// MessageSend implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) MessageSend(ctx context.Context, req *society.MessageSendRequest) (resp *society.MessageSendResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
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
	err = query.Q.MessageChat.Create(m)
	if err != nil {
		log.Printf("message send failed FromUserId:%d ToUserId:%d MsgContent:%s", myUserId, friendUserId, content)
		return nil, err
	}
	//cache
	result, _ := cache.Client.Exists(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId)).Result()
	if result > 0 {
		createTimeFormat := m.CreatedAt.Format("2006-01-02 15:04")
		//存在
		sM := &society.Message{
			Id:               int64(m.ID),
			FromUserId:       int64(m.FromUserId),
			ToUserId:         int64(m.ToUserId),
			Content:          m.MsgContent,
			CreateTime:       m.CreatedAt.UnixMilli(),
			CreateTimeFormat: &createTimeFormat,
		}
		bytes, _ := json.Marshal(sM)
		cache.Client.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), string(bytes))
		cache.Client.Expire(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), time.Hour)
	}
	result, _ = cache.Client.Exists(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", friendUserId, myUserId)).Result()
	if result > 0 {
		createTimeFormat := m.CreatedAt.Format("2006-01-02 15:04")
		//存在
		sM := &society.Message{
			Id:               int64(m.ID),
			FromUserId:       int64(m.FromUserId),
			ToUserId:         int64(m.ToUserId),
			Content:          m.MsgContent,
			CreateTime:       m.CreatedAt.UnixMilli(),
			CreateTimeFormat: &createTimeFormat,
		}
		bytes, _ := json.Marshal(sM)
		cache.Client.RPush(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), string(bytes))
		cache.Client.Expire(ctx, fmt.Sprintf("society.rpc:messagechat:userId(%d|%d)", myUserId, friendUserId), time.Hour)
	}
	return &society.MessageSendResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// IsFriend implements the SocietyServiceImpl interface.
func (s *SocietyServiceImpl) IsFriend(ctx context.Context, req *society.IsFriendRequest) (resp *society.IsFriendResponse, err error) {
	err = checkReq(req)
	if err != nil {
		return nil, err
	}
	r1, err := query.Q.Relation.Where(query.Q.Relation.FromUserId.Eq(uint(req.MyUserId)), query.Q.Relation.ToUserId.Eq(uint(req.FriendUserId)), query.Q.Relation.RelType.Eq(1)).First()
	if err != nil {
		klog.Infof("%d %d not friend", req.MyUserId, req.FriendUserId)
		return nil, errors.New(fmt.Sprintf("%d %d not friend", req.MyUserId, req.FriendUserId))
	}
	r2, err := query.Q.Relation.Where(query.Q.Relation.FromUserId.Eq(uint(req.FriendUserId)), query.Q.Relation.ToUserId.Eq(uint(req.MyUserId)), query.Q.Relation.RelType.Eq(1)).First()
	if err != nil {
		klog.Infof("%d %d not friend", req.MyUserId, req.FriendUserId)
		return nil, errors.New(fmt.Sprintf("%d %d not friend", req.MyUserId, req.FriendUserId))
	}
	if r1 != nil && r2 != nil {
		return &society.IsFriendResponse{StatusCode: 0, StatusMsg: "success"}, nil
	}
	return nil, errors.New(fmt.Sprintf("%d %d not friend", req.MyUserId, req.FriendUserId))
}

// 校验req是否为空，空值直接返回err
func checkReq(req interface{}) error {
	if req == nil {
		klog.Warnf("req is nil please check other service")
		return errors.New("req is nil please check other service")
	}
	return nil
}
