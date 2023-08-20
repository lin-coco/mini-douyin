package society_rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"mini-douyin/internal/basics_rpc"
	"mini-douyin/internal/pkg/dal/query"
	"mini-douyin/internal/pkg/kitex_gen/douyin/basics/basicsservice"
	EtcdConfig "mini-douyin/pkg/etcd"
	"mini-douyin/pkg/rdb"
	"sync"
	"time"
)

var (
	Config          *configuration
	BasicsRpcClient basicsservice.Client
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		Config = new(configuration).Init()
		BasicsRpcClient = BasicsRPCInit()
		query.SetDefault(rdb.DB)
	})
}

func BasicsRPCInit() basicsservice.Client {
	r, err := etcd.NewEtcdResolver(EtcdConfig.Config.Etcd.EndPoints)
	if err != nil {
		panic(err)
	}
	c, err := basicsservice.NewClient(
		basics_rpc.Config.ServerName,
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r))                            // resolver
	if err != nil {
		panic(err)
	}
	return c
}
