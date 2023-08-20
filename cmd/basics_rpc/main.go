package main

import (
	"fmt"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"mini-douyin/cmd/pkg"
	"mini-douyin/internal/basics_rpc"
	"mini-douyin/internal/pkg/kitex_gen/douyin/basics/basicsservice"
	EtcdConfig "mini-douyin/pkg/etcd"
)

func main() {
	//注册etcd
	registry, err := etcd.NewEtcdRegistry(EtcdConfig.Config.Etcd.EndPoints)
	if err != nil {
		log.Fatal(err)
	}
	svr := basicsservice.NewServer(new(basics_rpc.BasicsServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: basics_rpc.Config.ServerName}),
		server.WithRegistry(registry),
		server.WithServiceAddr(pkg.Addr{Net: "tcp", Address: fmt.Sprintf("0.0.0.0:%d", basics_rpc.Config.Port)}),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
