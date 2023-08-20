package cache

import (
	"fmt"
	"github.com/spf13/viper"
)

type configuration struct {
	Cache struct {
		Addr     string
		Password string
		PoolSize int
	} `json:"cache"`
}

func (c *configuration) Init() *configuration {
	return c.Load().Defaults().Fatal()
}

func (c *configuration) Load() *configuration {
	v := viper.New()
	v.AddConfigPath("../../configs")
	v.SetConfigName("cache.yaml")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	if err := v.Unmarshal(c); err != nil {
		panic(err)
	}
	return c
}

func (c *configuration) Defaults() *configuration {
	return c
}

func (c *configuration) Fatal() *configuration {
	return c
}
