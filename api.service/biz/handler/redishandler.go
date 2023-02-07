package handler

import (
	"api.service/biz/model/api/douyin/core"
	"api.service/biz/redis"
	"context"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"time"
)

//RedisSet 新增
//key string
//value map[int64]map[string]interface{}

type rMap map[string]interface{}
type rList []interface{}

func VideoRedisSet(videos []*core.Video) {
	//var err error
	//ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	//defer cancel()

	//err = redis.RedisDB.Set(ctx, "videos", , time.Hour).Err()
	//if err != nil {
	//	log.Infof("redis set videos")
	//}
}

func UserRedisSet(users []*core.User) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	m := make(map[int64]map[string]interface{}, 0)
	err = redis.RedisDB.Set(ctx, "videos", m, time.Hour).Err()
	if err != nil {
		log.Infof("redis set videos")
	}
}
