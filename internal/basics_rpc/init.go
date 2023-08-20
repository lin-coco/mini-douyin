package basics_rpc

import (
	"mini-douyin/internal/pkg/dal/query"
	"mini-douyin/pkg/rdb"
	"sync"
)

var (
	Config *configuration
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		Config = new(configuration).Init()
		query.SetDefault(rdb.DB)
	})
}
