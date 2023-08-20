package cache

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	Config *configuration
	Client *redis.Client
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		Config = new(configuration).Init()
		redisDB := redis.NewClient(&redis.Options{
			Addr:     Config.Cache.Addr,
			Password: Config.Cache.Password,
			PoolSize: Config.Cache.PoolSize,
		})
		Client = redisDB
	})
}
