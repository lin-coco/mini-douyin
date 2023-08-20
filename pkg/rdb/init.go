package rdb

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"sync"
	"time"
)

var (
	Config *configuration
	DB     *gorm.DB
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		Config = new(configuration).Init()
		fmt.Println(Config)
		DB = ConnectDB()
	})
}

func ConnectDB() *gorm.DB {
	if Config.Database.Separation {
		return ConnectRWDB()
	}
	dsn := Config.Database.Master
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db, err := d.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)
	return d
}

func ConnectRWDB() *gorm.DB {
	logrus.Info("使用读写分离")
	dsn := Config.Database.Master
	d, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}))
	if err != nil {
		panic(err)
	}
	var replicas []gorm.Dialector
	for i, s := range Config.Database.Slave {
		cfg := mysql.Config{
			DSN: s,
		}
		logrus.Infof("读写分离-%d-%s", i, s)
		replicas = append(replicas, mysql.New(cfg))
	}

	err = d.Use(
		dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: dsn,
			})},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(10).
			SetConnMaxLifetime(time.Hour).
			SetMaxOpenConns(200),
	)
	if err != nil {
		panic(err)
	}
	return d
}
