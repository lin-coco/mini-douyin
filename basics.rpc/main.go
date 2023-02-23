package main

import (
	"basics.rpc/dal/query"
	core "basics.rpc/kitex_gen/douyin/core/basicsservice"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
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
	OSSInit()
	InitRedis()
	//LogInit()
	//注册etcd
	//registry, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2479", "127.0.0.1:2579", "127.0.0.1:2679"})
	registry, err := etcd.NewEtcdRegistry([]string{"101.132.182.230:2379"})
	if err != nil {
		log.Fatal(err)
	}
	svr := core.NewServer(new(BasicsServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "basics.rpc"}),
		server.WithRegistry(registry),
		server.WithServiceAddr(addr{"tcp", "127.0.0.1:8890"}),
		//server.WithServiceAddr(addr{"tcp", "127.0.0.1:9890"}),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// 数据库
var DB *gorm.DB

var RedisDB *redis.Client

// oss实例
var OSSClient *oss.Client
var OSSBucket *oss.Bucket

const OSSBaseUrl string = "https://mini-douyin-basics.oss-cn-nanjing.aliyuncs.com/"
const BucketName = "mini-douyin-basics"

// 密码盐值
const PwdKey = "xys_"

func LogInit() {
	f, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	klog.SetOutput(f)
}

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
}

func OSSInit() {
	endpoint := "oss-cn-nanjing.aliyuncs.com"
	accessKey := "LTAI5tEZhWv2jdYrUKRQjm6Q"
	accessSecret := "1eqXnptJSUBpQxLirDYxsNqBfiHgir"
	bucketName := "mini-douyin-basics"
	//创建oss实例
	client, err := oss.New(endpoint, accessKey, accessSecret)
	if err != nil {
		panic("failed to connect oss")
	}
	OSSClient = client
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		panic("failed to connect bucket")
	}
	OSSBucket = bucket
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
