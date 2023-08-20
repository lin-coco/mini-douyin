package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

var (
	Config     *configuration
	EtcdClient *clientv3.Client
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		Config = new(configuration).Init()
		EtcdClient = initEtcdClient()
	})
}

func initEtcdClient() *clientv3.Client {

	// 配置
	config := clientv3.Config{
		Endpoints:   []string{"http://10.211.55.11:2379"},
		DialTimeout: 3 * time.Second,
	}

	// 连接
	c, err := clientv3.New(config)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
