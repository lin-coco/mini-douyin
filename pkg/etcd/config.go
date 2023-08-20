package etcd

import (
	"github.com/spf13/viper"
)

type configuration struct {
	Etcd struct {
		EndPoints   []string
		DialTimeout int
	} `json:"etcd"`
}

func (c *configuration) Init() *configuration {
	return c.Load().Defaults().Fatal()
}

func (c *configuration) Load() *configuration {
	vip := viper.New()
	vip.AddConfigPath("../../configs") //设置读取的文件路径
	vip.SetConfigName("etcd")          //设置读取的文件名
	vip.SetConfigType("yaml")          //设置文件的类型
	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}
	err := vip.Unmarshal(c)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *configuration) Defaults() *configuration {
	if c.Etcd.DialTimeout == 0 {
		c.Etcd.DialTimeout = 10
	}
	return c
}

func (c *configuration) Fatal() *configuration {
	if len(c.Etcd.EndPoints) == 0 {
		panic("end-points fatal")
	}
	if c.Etcd.DialTimeout < 0 {
		panic("dial-timeout fatal")
	}
	return c
}
