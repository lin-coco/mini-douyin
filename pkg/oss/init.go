package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"sync"
)

// oss实例
var (
	Config    *configuration
	OSSClient *oss.Client
	OSSBucket *oss.Bucket
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		Config = new(configuration).Init()
		endpoint := Config.OSS.Endpoint
		accessKey := Config.OSS.AccessKey
		accessSecret := Config.OSS.AccessSecret
		bucketName := Config.OSS.BucketName
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
	})
}
