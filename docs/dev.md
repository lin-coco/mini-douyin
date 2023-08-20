# mini-douyin 2.0.0 开发文档

## 基础架构

- 网关层
  - api.service：负责权限校验、安全可靠、限流熔断、负载均衡
- 业务层
  - basics.rpc：基础服务，负责基础业务的功能
  - interaction.rpc：互动服务，负责互动业务的功能
  - society.rpc：社交服务，负责社交业务的功能

流量必须先经过网关层再经过业务层

- 存储层
  - MySQL：负责用户数据存储
  - OSS：负责文件存储
  - ElasticSearch（待定）：负责日志等非用户数据
- 缓存层
  - Redis：负责业务层的性能优化

## 技术选型

1. Hertz：http框架，适用于实现网关
2. Coraza：WAF引擎，适用于实现网关层安全可靠
3. Kitex：rpc框架，适用于业务层框架
4. Gorm：orm框架，适用于操作mysql数据库
5. Etcd：注册中心、配置中心、KV存储，适用于注册服务，存储关键配置
6. MySQL：关系型数据库，存储业务层数据
7. OSS：对象存储，存储用户视频图片等文件
8. ElasticSearch：搜索与分析引擎，存储网关与业务层日志等数据
9. Redis：KV缓存，适用于业务层的缓存

## 代码结构

遵循`golang-standards`规范

```
├── api
├── assets
├── build
│   ├── ci
│   └── package
├── cmd
│   └── _your_app_
├── configs
├── deployments
├── docs
├── examples
├── githooks
├── init
├── internal
│   ├── app
│   │   └── _your_app_
│   └── pkg
│       └── _your_private_lib_
├── pkg
│   └── _your_public_lib_
├── scripts
├── test
├── third_party
├── tools
├── vendor
├── web
│   ├── app
│   ├── static
│   └── template
├── website
├── .gitignore
├── LICENSE.md
├── Makefile
├── README.md
└── go.mod
```

## 环境搭建

（Linux版本）

### go

下载：`wget https://go.dev/dl/go1.19.11.linux-arm64.tar.gz`

解压：`tar -C /usr/local -xzf go1.19.11.linux-arm64.tar.gz`

添加环境变量：`echo "export PATH=/usr/local/go/bin:${PATH}" >> /etc/environment && source /etc/environment`

检查版本：`go version`

修改GOPROXY：`go env -w GOPROXY=https://goproxy.cn,direct`

### etcd集群

执行脚本：

```sh
ETCD_VER=v3.4.27

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}

rm -f /tmp/etcd-${ETCD_VER}-linux-arm64.tar.gz
rm -rf /tmp/etcd && mkdir -p /tmp/etcd

curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-arm64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-arm64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-arm64.tar.gz -C /tmp/etcd --strip-components=1
rm -f /tmp/etcd-${ETCD_VER}-linux-arm64.tar.gz

/tmp/etcd/etcd --version
/tmp/etcd/etcdctl version

cp /tmp/etcd/etcd /usr/local/bin/etcd
cp /tmp/etcd/etcd /usr/local/bin/etcdctl

echo "export ETCD_UNSUPPORTED_ARCH=arm64" >> /etc/environment
echo "export ETCDCTL_API=3" >> /etc/environment
source /etc/environment
```

加载环境变量：`source /etc/enbironment`



静态集群搭建

添加环境变量

```sh
echo "export ETCD_INITIAL_CLUSTER="infra0=http://10.0.1.10:2380,infra1=http://10.0.1.11:2380,infra2=http://10.0.1.12:2380" >> /etc/environment
echo "export ETCD_INITIAL_CLUSTER_STATE=new" >> /etc/enbironment
```

在每个etcd服务上执行

```sh
$ etcd --name infra0 --initial-advertise-peer-urls http://10.211.55.11:2380 \
  --listen-peer-urls http://10.211.55.11:2380 \
  --listen-client-urls http://10.211.55.11:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.211.55.11:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://10.211.55.11:2380,infra1=http://10.211.55.12:2380,infra2=http://10.211.55.13:2380 \
  --initial-cluster-state new
```

```sh
$ etcd --name infra1 --initial-advertise-peer-urls http://10.211.55.12:2380 \
  --listen-peer-urls http://10.211.55.12:2380 \
  --listen-client-urls http://10.211.55.12:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.211.55.12:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://10.0.1.10:2380,infra1=http://10.0.1.11:2380,infra2=http://10.0.1.12:2380 \
  --initial-cluster-state new
```

```sh
$ etcd --name infra2 --initial-advertise-peer-urls http://10.211.55.13:2380 \
  --listen-peer-urls http://10.211.55.13:2380 \
  --listen-client-urls http://10.211.55.13:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.211.55.13:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://10.211.55.11:2380,infra1=http://10.211.55.12:2380,infra2=http://10.211.55.13:2380 \
  --initial-cluster-state new
```

### mysql主从



### redis主从



## 技术复习

### Hertz

略

### Kitex

略

### GORM

略



...



## Coraza

### Bot防护

https://www.anquanke.com/post/id/273031

主要功能

1. Bot特征检测和识别：

   系统基于已知的Bot特征标签库

   - 初步划分出人为正常流量、Bot合法流量、Bot恶意流量与未知流量

   - 按照危险等级划分：低、中、高危

   - Bot机器人类型细分：间谍程序类Bot、搜索引擎类Bot与拒绝服务类Bot等20大类

   - 根据客户端Bot类型，分别进行精准访问控制和限速等其他缓解策略

2.  基于场景的分级筛选模型

   - 信息收集
   - 流量交互
   - 行为分析
   - 差分计算
   - 风险流量过滤
   - 恶意流量分解
   - 动态业务场景防御策略模型



### API安全



## 造轮子系列

### 分布式唯一id生成器

#### 要求

全局不重复 + 不可猜测 + 递增态势

#### 设计思路

基于Mist薄雾算法

位数：64位

占位：

- 第一段为最高位，占 1 位，保持为 0，使得值永远为正数；
- 第二段放置自增数，占 47 位，自增数在高位能保证结果值呈递增态势，遂低位可以为所欲为；
- 第三段放置随机因子一，占 8 位，上限数值 255，使结果值不可预测；
- 第四段放置随机因子二，占 8 位，上限数值 255，使结果值不可预测；

薄雾算法优点

1. 不依赖时钟，效率高
2. 全局不重复、不可猜测、递增态势

薄雾算法缺点

1. 重新启动会造成递增数值回到初始值



UidGenerator优点

1. RingBuffer缓存
2. 借用未来时间解决时钟回拨

弥补薄雾算法重启递增数值回到初始值



占位：

- 第一段为最高位，占 1 位，保持为 0，使得值永远为正数；
- 第二段放置时间戳，占 41 位；
- 第三段放置自增数，占 14 位；
- 第四段放置随机因子，占 8 位，上限为255，使结果值不可预测；

UIDRingBuffer：缓存生成的id

FlagRingBuffer：记录uid状态（是否可填充、是否可消费）

CurUsedTimestamp：当前uid计算时间戳

RandChannel：生成随机因子





BUID优势

1. 全局不重复
2. 局部不可猜测
3. 递增态势

4. 缓存











