package main

import (
	"basics.rpc/kitex_gen/douyin/core/basicsservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"society.rpc/dal/query"
	second "society.rpc/kitex_gen/douyin/extra/second/societyservice"
	"time"
)

type addr struct {
	network string
	address string
}

func (a addr) Network() string {
	return a.network
}
func (a addr) String() string {
	return a.address
}

func main() {
	//初始化其他服务
	DBInit()
	BasicsRPCInit()
	InitRedis()
	//注册etcd
	registry, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	svr := second.NewServer(new(SocietyServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "society.rpc"}),
		server.WithRegistry(registry),
		server.WithServiceAddr(addr{"tcp", "127.0.0.1:8892"}),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// 数据库
var DB *gorm.DB
var Q *query.Query
var RedisDB *redis.Client

// basics.rpc
var BasicsService basicsservice.Client

func DBInit() {
	dsn := "xys:232020ctt@@tcp(rm-uf6e4xr978w748b9w7o.mysql.rds.aliyuncs.com:3306)/sql_test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{ //连接数据库
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		AllowGlobalUpdate:      true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	query.SetDefault(db)
	Q = query.Q
}

func BasicsRPCInit() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	c, err := basicsservice.NewClient(
		"basics.rpc",
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r))                            // resolver
	if err != nil {
		panic(err)
	}
	BasicsService = c
}

func InitRedis() {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "r-uf6tmk8szjtlbvorcypd.redis.rds.aliyuncs.com:6379",
		Password: "syr1120@xyscom",
		DB:       0,
		PoolSize: 20,
	})
	RedisDB = redisDB
}
