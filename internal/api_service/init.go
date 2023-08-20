package api_service

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"mini-douyin/internal/basics_rpc"
	"mini-douyin/internal/interaction_rpc"
	"mini-douyin/internal/pkg/kitex_gen/douyin/basics/basicsservice"
	"mini-douyin/internal/pkg/kitex_gen/douyin/interaction/interactionservice"
	"mini-douyin/internal/pkg/kitex_gen/douyin/society/societyservice"
	"mini-douyin/internal/society_rpc"
	EtcdConfig "mini-douyin/pkg/etcd"
	"sync"
	"time"
)

var (
	Config               *configuration
	BasicsRpcClient      basicsservice.Client
	InteractionRpcClient interactionservice.Client
	SocietyRpcClient     societyservice.Client
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		Config = new(configuration).Init()
		BasicsRpcClient = initBasicsServiceRpc()
		InteractionRpcClient = initInteractionServiceRpc()
		SocietyRpcClient = initSocietyServiceRPC()
		InitVD()
	})
}

func initBasicsServiceRpc() basicsservice.Client {
	r, err := etcd.NewEtcdResolver(EtcdConfig.Config.Etcd.EndPoints)
	if err != nil {
		panic(err)
	}
	c, err := basicsservice.NewClient(
		basics_rpc.Config.ServerName,
		client.WithMuxConnection(1),                                // mux
		client.WithRPCTimeout(3*time.Second),                       // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),             // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()),          // retry
		client.WithSuite(trace.NewDefaultClientSuite()),            // tracer
		client.WithResolver(r),                                     // resolver
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer())) // 负载均衡)

	if err != nil {
		panic(err)
	}
	return c
}

func initInteractionServiceRpc() interactionservice.Client {
	r, err := etcd.NewEtcdResolver(EtcdConfig.Config.Etcd.EndPoints)
	if err != nil {
		panic(err)
	}
	c, err := interactionservice.NewClient(
		interaction_rpc.Config.ServerName,
		client.WithMuxConnection(1),                                // mux
		client.WithRPCTimeout(3*time.Second),                       // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),             // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()),          // retry
		client.WithSuite(trace.NewDefaultClientSuite()),            // tracer
		client.WithResolver(r),                                     // resolver
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer())) // 负载均衡)
	if err != nil {
		panic(err)
	}
	return c
}

func initSocietyServiceRPC() societyservice.Client {
	r, err := etcd.NewEtcdResolver(EtcdConfig.Config.Etcd.EndPoints)
	if err != nil {
		panic(err)
	}
	c, err := societyservice.NewClient(
		society_rpc.Config.ServerName,
		client.WithMuxConnection(1),                                // mux
		client.WithRPCTimeout(3*time.Second),                       // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),             // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()),          // retry
		client.WithSuite(trace.NewDefaultClientSuite()),            // tracer
		client.WithResolver(r),                                     // resolver
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer())) // 负载均衡)
	if err != nil {
		panic(err)
	}
	return c
}
