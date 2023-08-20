package main

import (
	"fmt"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"mini-douyin/cmd/pkg"
	"mini-douyin/internal/interaction_rpc"
	"mini-douyin/internal/pkg/kitex_gen/douyin/interaction/interactionservice"
	EtcdConfig "mini-douyin/pkg/etcd"
)

func main() {
	//注册etcd
	registry, err := etcd.NewEtcdRegistry(EtcdConfig.Config.Etcd.EndPoints)
	if err != nil {
		log.Fatal(err)
	}

	svr := interactionservice.NewServer(
		new(interaction_rpc.InteractionServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: interaction_rpc.Config.ServerName}),
		server.WithRegistry(registry),
		server.WithServiceAddr(pkg.Addr{Net: "tcp", Address: fmt.Sprintf("0.0.0.0:%d", interaction_rpc.Config.Port)}),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
