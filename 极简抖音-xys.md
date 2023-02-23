# 极简抖音

## 预开发

### 项目设计要求

一个基础功能两大方向：基础功能、互动方向、社交方向

1. 基础功能
   - 视频Feed流：支持所有用户刷抖音，视频按照投稿时间倒序输出
   - 视频投稿：支持登录用户自己拍视频投稿
   - 个人主页：支持查看用户基本信息和投稿列表，注册用户流程简化
2. 互动方向
   - 喜欢列表：登录用户可以对视频点赞，在个人主页喜欢Tab下能够查看点赞视频列表
   - 评论列表：支持未登录用户查看视频下的评论列表，登录用户能够发表评论
3. 社交方向
   - 登录用户可以关注其他用户，能够在个人主页查看本人的关注数和粉丝数，查看关注列表和粉丝列表
   - 登录用户在消息页展示已关注的用户列表，点击用户头像进入聊天后可以发送消息

### 微服务方案

目前了解有两种方案

1. go-zero
2. hertz + kitex + gorm

都是非常简单实用、特性丰富的框架。开发模式基本差不多都可以选择，本项目使用hertz

### 微服务拆分

![架构图的副本](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/%E6%9E%B6%E6%9E%84%E5%9B%BE%E7%9A%84%E5%89%AF%E6%9C%AC.png)

> api服务和各种rpc服务可以是单体或者集群，依据最终部署条件确定

1. 会话服务

   - 实现所有服务功能的api
   - api层：调用各个rpc实现业务逻辑，对外提供接口

2. 基础服务

   - 实现用户和视频的相关方法

   - 可能的表：用户、视频
   - rpc层：对内提供用户功能和视频功能的调用
   - 不需要调用其他rpc服务

3. 互动服务

   - 实现互动方向的相关方法
   - 可能的表：用户点赞表——关联用户id与点赞视频id；评论表——评论用户、评论信息……（被评论主体只有视频，不需要树形结构）
   - 需要调用**基础服务**的rpc，获得用户信息视频信息等

4. 社交服务

   - 实现社交方向的相关方法
   - 可能的表：关注与粉丝的表结构设计有多种，尽可能保证在搜索关注列表和粉丝列表都要高效，消息表
   - 需要调用**基础服务**的rpc，获得用户信息等

没有循环调用的rpc服务

如果所有服务本地都有代码，那么可以通过go module导入本地包

如果服务代码均在在不同人电脑上，那必须上传到github当中。关于rpc方法的调用源码可以通过go get获取 或者 直接复制代码到调用方

业务逻辑很简单，就不细说了

### 会话服务

api服务；提供所有功能的api，需要进行token认证；连接etcd，调用rpc，实现service部分的业务逻辑；

token说明：app自动重启会丢失token，登录过的用户会有一个token，用jwt生成，有效期不关心，测试时可以长一点，部署时可以短一点。token关键信息必须包含userId。

**基础接口**

```protobuf
syntax = "proto3";
package api.douyin.core;
option go_package="api/douyin/core";

import "api.proto";

message DouyinFeedRequest {
  optional int64 latest_time = 1[(api.query)="latest_time"]; // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
  optional string token = 2[(api.query)="token"]; // 可选参数，登录用户设置
}
message DouyinFeedResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  repeated Video video_list = 3; // 视频列表
  optional int64 next_time = 4; // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

message Video {
  int64 id = 1; // 视频唯一标识 18
  User author = 2; // 视频作者信息
  string play_url = 3; // 视频播放地址
  string cover_url = 4; // 视频封面地址
  int64 favorite_count = 5; // 视频的点赞总数
  int64 comment_count = 6; // 视频的评论总数
  bool is_favorite = 7; // true-已点赞，false-未点赞
  string title = 8; // 视频标题
}
message User {
  int64 id = 1; // 用户id
  string name = 2; // 用户名称
  optional int64 follow_count = 3; // 关注总数
  optional int64 follower_count = 4; // 粉丝总数
  bool is_follow = 5; // true-已关注，false-未关注
}

message DouyinUserRegisterRequest {
  string username = 1[(api.query)="username"]; //注册用户名，最长32个字符
  string password = 2[(api.query)="password"]; //密码，最长32个字符
}

message DouyinUserRegisterResponse {
  int32 status_code = 1; //状态码，0-成功，其他值-失败
  optional string status_msg = 2; //返回状态描述
  int64 user_id = 3; //用户id
  string token = 4; //用户鉴权token
}

message DouyinUserLoginRequest {
  string username = 1[(api.query)="username"]; // 登录用户名
  string password = 2[(api.query)="password"]; // 登录密码
}
message DouyinUserLoginResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  int64 user_id = 3; // 用户id
  string token = 4; // 用户鉴权token
}

message DouyinUserRequest {
  int64 user_id = 1[(api.query)="user_id"]; // 用户id
  string token = 2[(api.query)="token"]; // 用户鉴权token 7
}
message DouyinUserResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  User user = 3; // 用户信息
}

message DouyinPublishActionRequest {
  string token = 1[(api.body)="token"]; // 用户鉴权token
  bytes data = 2[(api.body)="data"]; // 视频数据
  string title = 3[(api.body)="title"]; // 视频标题
}
message DouyinPublishActionResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
}

message DouyinPublishListRequest {
  int64 user_id = 1[(api.query)="user_id"]; // 用户id
  string token = 2[(api.query)="token"]; // 用户鉴权token
}
message DouyinPublishListResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述 12
  repeated Video video_list = 3; // 用户发布的视频列表
}


service CoreService {
  // 视频流接口
  rpc FeedRequest(DouyinFeedRequest) returns(DouyinFeedResponse) {
    option (api.get) = "/douyin/feed";
  }
  // 用户注册接口
  rpc RegisterRequest(DouyinUserRegisterRequest) returns(DouyinUserRegisterResponse) {
    option (api.post) = "/douyin/user/register";
  }
  // 用户登录接口
  rpc LoginRequest(DouyinUserLoginRequest) returns(DouyinUserLoginResponse) {
    option (api.post) = "/douyin/user/login";
  }
  // 用户信息
  rpc UserRequest(DouyinUserRequest) returns(DouyinUserResponse) {
    option (api.get) = "/douyin/user";
  }
  // 视频投稿
  rpc PublishActionRequest(DouyinPublishActionRequest) returns(DouyinPublishActionResponse) {
    option (api.post) = "/douyin/publish/action";
  }
  // 发布列表
  rpc PublishListRequest(DouyinPublishListRequest) returns(DouyinPublishListResponse) {
    option (api.get) = "/douyin/publish/list";
  }
}
```

**互动接口**

```protobuf
syntax = "proto3";

package api.douyin.extra.first;
option go_package="api/douyin/extra/first";

import "api.proto";

message DouyinFavoriteActionRequest {
  string token = 1[(api.query)="token"]; //用户鉴权token
  int64 video_id = 2[(api.query)="video_id"]; //视频id
  int32 action_type = 3[(api.query)="action_type"]; //1-点赞，2-取消点赞
}

message DouyinFavoriteActionResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
}

message DouyinFavoriteListRequest {
  int64 user_id = 1[(api.query)="user_id"]; // 用户id
  string token = 2[(api.query)="token"]; // 用户鉴权token
}

message DouyinFavoriteListResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  repeated Video video_list = 3; // 用户点赞视频列表
}
message Video {
  int64 id = 1; // 视频唯一标识
  User author = 2; // 视频作者信息
  string play_url = 3; // 视频播放地址
  string cover_url = 4; // 视频封面地址
  int64 favorite_count = 5; // 视频的点赞总数
  int64 comment_count = 6; // 视频的评论总数
  bool is_favorite = 7; // true-已点赞，false-未点赞
  string title = 8; // 视频标题
}
message User {
  int64 id = 1; // 用户id
  string name = 2; // 用户名称 29
  optional int64 follow_count = 3; // 关注总数
  optional int64 follower_count = 4; // 粉丝总数
  bool is_follow = 5; // true-已关注，false-未关注
}

message DouyinCommentActionRequest {
  string token = 1[(api.query)="token"]; // 用户鉴权token
  int64 video_id = 2[(api.query)="video_id"]; // 视频id
  int32 action_type = 3[(api.query)="action_type"]; // 1-发布评论，2-删除评论
  optional string comment_text = 4[(api.query)="comment_text"]; // 用户填写的评论内容，在action_type=1的时候使用
  optional int64 comment_id = 5[(api.query)="comment_id"]; // 要删除的评论id，在action_type=2的时候使用
}
message DouyinCommentActionResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  optional Comment comment = 3; // 评论成功返回评论内容，不需要重新拉取整个列表
}
message Comment {
  int64 id = 1; // 视频评论id
  User user =2; // 评论用户信息
  string content = 3; // 评论内容
  string create_date = 4; // 评论发布日期，格式 mm-dd
}

message DouyinCommentListRequest {
  string token = 1[(api.query)="token"]; // 用户鉴权token
  int64 video_id = 2[(api.query)="video_id"]; // 视频id
}
message DouyinCommentListResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  repeated Comment comment_list = 3; // 评论列表
}

service InteractionService {
  //赞操作
  rpc FavoriteAction(DouyinFavoriteActionRequest) returns(DouyinFavoriteActionResponse) {
    option (api.post) = "/douyin/favorite/action";
  }
  //喜欢列表
  rpc FavoriteList(DouyinFavoriteListRequest) returns(DouyinFavoriteListResponse) {
    option (api.get) = "/douyin/favorite/list";
  }
  //评论操作
  rpc CommentAction(DouyinCommentActionRequest) returns(DouyinCommentActionResponse) {
    option (api.post) = "/douyin/comment/action";
  }
  //视频评论列表
  rpc CommentList(DouyinCommentListRequest) returns(DouyinCommentListResponse) {
    option (api.get) = "/douyin/comment/list";
  }
}
```

**社交接口**

```protobuf
syntax = "proto3";

package api.douyin.extra.second;
option go_package = "api/douyin/extra/second";

import "api.proto";

message DouyinRelationActionRequest {
  string token = 1[(api.query)="token"]; // 用户鉴权token
  int64 to_user_id = 2[(api.query)="to_user_id"]; // 对方用户id
  int32 action_type = 3[(api.query)="action_type"]; // 1-关注，2-取消关注
}
message DouyinRelationActionResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
}

message DouyinRelationFollowListRequest {
  int64 user_id = 1[(api.query)="user_id"]; // 用户id
  string token = 2[(api.query)="token"]; // 用户鉴权token
}
message DouyinRelationFollowListResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  repeated User user_list = 3; // 用户信息列表
}
message User {
  int64 id = 1; // 用户id
  string name = 2; // 用户名称
  optional int64 follow_count = 3; // 关注总数
  optional int64 follower_count = 4; // 粉丝总数
  bool is_follow = 5; // true-已关注，false-未关注
}

message DouyinRelationFollowerListRequest {
  int64 user_id = 1[(api.query)="user_id"]; // 用户id
  string token = 2[(api.query)="token"]; // 用户鉴权token
}
message DouyinRelationFollowerListResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  repeated User user_list = 3; // 用户列表
}

message DouyinRelationFriendListRequest {
  int64 user_id = 1[(api.query)="user_id"]; // 用户id
  string token = 2[(api.query)="token"]; // 用户鉴权token
}
message DouyinRelationFriendListResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  repeated User user_list = 3; // 用户列表
}

message DouyinMessageChatRequest {
  string token = 1[(api.query)="token"]; // 用户鉴权token
  int64 to_user_id = 2[(api.query)="to_user_id"]; // 对方用户id
}

message DouyinMessageChatResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
  repeated Message message_list = 3; // 消息列表
}

message Message {
  int64 id = 1; // 消息id
  int64 to_user_id = 2; // 该消息接收者的id
  int64 from_user_id =3; // 该消息发送者的id
  string content = 4; // 消息内容
  optional string create_time = 5; // 消息创建时间
}

message DouyinMessageActionRequest {
  string token = 1; // 用户鉴权token
  int64 to_user_id = 2; // 对方用户id
  int32 action_type = 3; // 1-发送消息
  string content = 4; // 消息内容
}

message DouyinMessageActionResponse {
  int32 status_code = 1; // 状态码，0-成功，其他值-失败
  optional string status_msg = 2; // 返回状态描述
}


service SocietyService {
  //关系操作
  rpc RelationAction(DouyinRelationActionRequest) returns(DouyinRelationActionResponse) {
    option (api.post) = "/douyin/relation/action";
  }
  //用户关注列表
  rpc RelationFollowList(DouyinRelationFollowListRequest) returns(DouyinRelationFollowListResponse) {
    option (api.get) = "/douyin/relation/follow/list";
  }
  //用户粉丝列表
  rpc RelationFollowerList(DouyinRelationFollowerListRequest) returns(DouyinRelationFollowerListResponse) {
    option (api.get) = "/douyin/relation/follower/list";
  }
  //用户好友列表
  rpc RelationFriendList(DouyinRelationFriendListRequest) returns(DouyinRelationFriendListResponse) {
    option (api.get) = "/douyin/relation/friend/list";
  }
  //消息方案一
  //...不考虑，好像读取云端消息记录的说明
  //消息方案二
  //...使用消息方案二
  //聊天记录
  rpc MessageChat(DouyinMessageChatRequest) returns(DouyinMessageChatResponse) {
    option (api.get) = "/douyin/message/chat/";
  }
  //发送消息
  rpc MessageAction(DouyinMessageActionRequest) returns(DouyinMessageActionResponse) {
    option (api.post) = "/douyin/message/action/";
  }
}
```

### 基础服务

**user 表**

```go
type User struct {
   gorm.Model
   Name     string `gorm:"size:256"`
   Password string `gorm:"size:256"`
}
```

**video 表**

```go
type Video struct {
   gorm.Model
   UserId         uint
   PlayUrl        string `gorm:"size:500"`
   CoverUrl       string `gorm:"size:500"`
   
   Title          string `gorm:"size:50"`
}
```

实现对user表和video表的增删改查，提供增删改查的rpc调用，注册进etcd当中

### 互动服务

**user_favourite 表**——用户点赞表（用户-视频 1对多）

点赞是一个实时性的操作

在写操作非常频繁的情况下，可以先缓存到内存中，异步定时更新到mysql，减少对mysql数据库的压力

读写操作暂时放在redis中，异步定时更新到mysql

```go
type UserFavourite struct {
	gorm.Model
	UserId  uint
	VideoId uint
	status  uint8 `gorm:"default:1"` //点赞 状态为1 取消赞状态为0
}
```

> 冷热数据物理存储分开存储的思路是对的，冷热数据读写特性不同（冷数据的读写比例高于热数据），分开储存之后可以采用不同的cache策略，冷数据因为更新少可以直接同步一份至redis这类NOSQL服务，业务层直接从redis读取，减少对mysqldb的压力；热数据因更新较频繁，可以根据用户id（或者说uin）hash到多台写服务，并先写至写服务器的本地缓存中，再异步定时批量更新至mysql，减少对mysql的写压力。

comment 表——评论表（只有对video主体对评论，没有对评论的评论）

```go
type Comment struct {
   gorm.Model
   content    string `gorm:"size:256"`
   FromUserId uint
   ToVideoId uint
}
```

### 社交服务

relation 表——关注粉丝表

```go
type Relation struct {
   gorm.Model
   FromUserId uint
   ToUserId   uint
   // FromUserId 关注了 ToUserId；
   // 查询关注列表 即 select to_user_id from relation where from_user_id = ? and rel_type 
   // 查询粉丝列表 即 select from_user_id from relation where to_user_id = ?
   RelType uint8 `gorm:"default:1"` //1为有效 0为无效
}
```

> 查询朋友不包括自己

message_chat 表——聊天表

```go
type MessageChat struct {
   gorm.Model
   // FromUserId 给 ToUserId 发送的 MsgContent
   MsgContent string `gorm:size:256`
   FromUserId uint
   ToUserId uint
}
```

## 开发中

### 开发过程中遇到的困难

1. 本地导包的问题，我是把所有的服务放在一个文件夹里，使用go module本地导包来访问服务的

2. app网络错误，307错误，解决了好久，没有解决，但却让我了解了抓包软件的基本使用，307错误也是通过抓包知道的。不知道怎么会出现这个错误

3. 上传文件如何绑定文件对象，获取文件流。对req属性的绑定，以及pb文件设计的值位置有了更深入的认识，最后没有用c.BindAndValidate(&req)，而是自己调试找到字段位置绑定req

   但之后通过官方文档发现，绑定文件需要特定的结构

   ![image-20230203190057429](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20230203190057429.png)

4. 视频封面搞了好久没弄好，没有弄了

5. 遇到最多的问题是空指针，经验不足，经常吃亏，req为空指针是我最烦的

......

### 服务规范方面

1. 鉴权token放在中间件里验证

2. 参数校验中间件验证、handler验证

3. log规范

4. 统一异常处理规范

5. rpc req参数nil值处理

......

## 服务性能与安全可靠

### mysql优化

```go
user, err := query.Q.User.Where(query.User.ID.Eq(uint(userId))).First()
	videos, err := query.Q.Video.Order(query.Video.CreatedAt).Where(query.Video.CreatedAt.Lt(latest)).Limit(30).Find()
		user, err := query.Q.User.Where(query.User.ID.Eq(videos[i].UserId)).First()
	err = query.Q.Video.Create(&model.Video{
	videos, err := query.Q.Video.Where(query.Video.UserId.Eq(uint(userId))).Find()
		user, err := query.Q.User.Where(query.User.ID.Eq(uint(userId))).First()
	user, _ := query.Q.User.Where(query.User.Name.Eq(username)).First()
		video, err := query.Q.Video.Where(query.Video.ID.Eq(uint(videoId))).First()
		user, err := query.Q.User.Where(query.User.ID.Eq(video.UserId)).First()
	users, err := query.Q.User.Where(query.User.ID.In(userIdList2...)).Find()
```

- video createAt 索引
- video userId 索引
- user username 索引

```
userFavourite, _ := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).First()
	userFavourite, _ := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId))).First()
	userFavourites, err := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.Status.Eq(1)).Find()
		favoriteCount, err := Q.UserFavourite.Where(query.UserFavourite.VideoId.Eq(uint(video.Id)), query.UserFavourite.Status.Eq(1)).Count()
		commentCount, err := Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(video.Id))).Count()
	comment, err := Q.Comment.Where(query.Comment.ID.Eq(uint(commentId))).Select(query.Comment.ToVideoId).First()
	comments, err := Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(videoId))).Order(query.Comment.CreatedAt).Find()
	count, err = Q.UserFavourite.Where(query.UserFavourite.VideoId.Eq(uint(videoId)), query.UserFavourite.Status.Eq(1)).Count()
	count, err = Q.Comment.Where(query.Comment.ToVideoId.Eq(uint(videoId))).Count()
	favourite, _ := Q.UserFavourite.Where(query.UserFavourite.UserId.Eq(uint(userId)), query.UserFavourite.VideoId.Eq(uint(videoId)), query.UserFavourite.Status.Eq(1)).First()
	comment, err := Q.Comment.Where(query.Comment.ID.Eq(uint(commentId))).First()
```

- user_favorite user_id和video_id联合索引
- comment to_video_id和createdAt联合索引
- user_favorite video_id索引

```
relation, _ := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).First()
	relation, _ := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(fromUserId)), query.Relation.ToUserId.Eq(uint(toUserId))).First()
	relations, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
		followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
	relations, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
		followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
		followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(user.Id)), query.Relation.RelType.Eq(1)).Count()
	relations, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
	friends, err := Q.Relation.Where(query.Relation.FromUserId.In(toUserIds...), query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Find()
		followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(friend.Id)), query.Relation.RelType.Eq(1)).Count()
		followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(friend.Id)), query.Relation.RelType.Eq(1)).Count()
	followCount, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	followerCount, err := Q.Relation.Where(query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	count, err := Q.Relation.Where(query.Relation.FromUserId.Eq(uint(myId)), query.Relation.ToUserId.Eq(uint(userId)), query.Relation.RelType.Eq(1)).Count()
	messageChats, err := Q.MessageChat.Where(query.MessageChat.FromUserId.In(uint(myUserId), uint(friendUserId)), query.MessageChat.ToUserId.In(uint(myUserId), uint(friendUserId))).Order(query.MessageChat.CreatedAt).Find()
```

- relation from_user_id索引
- relation to_user_id索引
- message_chat from_user_id to_user_id createdAt索引

### Redis缓存优化

技术方案：

Go语言操作Redis的方案有以下几种：

1. 使用Redigo库：这是一个Go语言实现的Redis客户端库，具有简单易用的特点。
2. 使用Go-Redis库：这是一个基于Go语言的高效、简单易用的Redis客户端库。
3. 直接使用Redis协议：Go语言支持通过直接发送Redis协议命令进行Redis操作。

不同的方案适用于不同的场景，具体的选择取决于具体的需求。建议根据自己的需求，选择适合自己的方案。

我使用go-redis库

> 一开始使用的redis8版本，遇到一个HSet的问题，field字段一直是map的序列化json，查了好多资料没有解决，就很迷，后来用redis9版本才解决

假设app无法做缓存的情况下，那么服务端Redis的任务也会很重。有些缓存适合在app当中做，如用户个人信息、视频等，当然也可以在服务端Redis当中做，只是前者更好。

在项目中使用Redis缓存优化是能很大程度上提升性能，但是网上没有成熟的go http框架使用Redis的规范案例，在这里自己研究一下，如何在项目中使用Redis。

使用Redis的优势：

1. 性能优化

使用Redis的弊端：

1. 需要维护的代码变多
2. 可以造成用户端数据的延迟，严重的数据不一致
3. Redis带来的一系列问题，如宕机、缓存穿透、雪崩等

可以在两个层面对请求进行优化

1. api层面的优化
   - api.service 关注列表、粉丝列表、好友列表的优化
2. mysql层面的优化
   - basics.rpc 登录注册的优化
   - interaction.rpc 获取点赞总数的优化、获取评论总数的优化
   - society.rpc 消息的优化

api.service 关注列表、粉丝列表、好友列表的优化

1. 缓存关注列表、粉丝列表和好友列表的数据
2. 修改关注列表在自己进行关注/取消操作的时候，进行添加/删除某个元素，添加删除失败删除key；
3. 删除粉丝列表在其他用户进行关注/取消自己的时候，删除key；
4. 好友列表在自己进行关注操作的时候，进行删除

### 接口优化

1. 一些大V的粉丝列表可能太多，我们需要限制返回数
2. 一些人的聊天记录非常多，存储在本地是一种方式，但要是存储在服务端，每次从服务端拿聊天记录，可以每次返回最近50条聊天记录

......

### 服务性能优化

- 集群，负载均衡
- redis缓存优化
- mysql索引
- ...

### 安全可靠

- 严格按照token的用户信息来判断请求是否合理
- sql注入，gorm类似的框架都有考虑这种情况
- 用户密码进行加密存储
- 统一响应处理，统一异常处理，已经给出了统一响应格式
- ...

### 系统稳定运行

- 日志收集
- 限流、跟踪等
- ...

多数利用框架的一些特性来解决问题

## 部署方案

微服务部署方案有成熟的kubernetes，我了解的不多，kubernetes对云服务器要求较高，不予考虑了

部署成容器或者直接在几个服务器上运行服务，rpc服务部署成集群的话修改端口号就可以部署成多个服务，框架可以进行负载均衡。api服务部署成集群，可以用nginx帮助负载均衡一下







