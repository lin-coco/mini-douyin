package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"interaction.rpc/kitex_gen/douyin/extra/first/interactionservice"
	"time"
)

var InteractionService interactionservice.Client

func initInteractionServiceRpc() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2479", "127.0.0.1:2579", "127.0.0.1:2679"})
	if err != nil {
		panic(err)
	}
	c, err := interactionservice.NewClient(
		"interaction.rpc",
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
	InteractionService = c
}
