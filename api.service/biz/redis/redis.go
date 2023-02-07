package redis

import "github.com/redis/go-redis/v9"

var RedisDB *redis.Client

func InitRedis() {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "r-uf6tmk8szjtlbvorcypd.redis.rds.aliyuncs.com:6379",
		Password: "syr1120@xyscom",
		DB:       0,
		PoolSize: 20,
	})
	RedisDB = redisDB
}
