package main

import (
	"flag"
	"goshop/configs"
	"goshop/models"
	"goshop/pkg/mysql"
	"goshop/pkg/redis"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type CartServer struct {
	service.Service
}

var (
	server *CartServer
	db     *gorm.DB
	rdb    redis.IRdb
	once   sync.Once
)

func CartServerGetInstance() *CartServer {
	once.Do(func() {
		server = &CartServer{}
		server.Derived = server
	})
	return server
}

func (s *CartServer) Init() bool {
	var err error
	if !configs.ParseConfig() {
		glog.Errorln("[CartServer] parse config error.")
		return false
	}
	// RPC server
	if !rpcServerStart() {
		glog.Errorln("[CartServer] rpc server start error.")
		return false
	}

	// MySQL connect
	if db, err = mysql.DatabaseInit(&configs.GetConf().MysqlCfg); err != nil {
		glog.Errorln("[CartServer] mysql database init error.")
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.Cart{})

	// Redis connect
	if rdb, err = redis.NewRedisClient(&configs.GetConf().RedisCfg); err != nil {
		glog.Errorln("[CartServer] redis database init error.")
		return false
	}
	glog.Infof("[CartServer] redis database [%s] init success.\n", configs.GetConf().GetRedisAddr())

	// Consul register
	if !service.ServiceRegister(
		"1",
		"cart-service",
		configs.GetConf().CartCfg.Host,
		configs.GetConf().CartCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[CartServer] consul register error.")
		return false
	}
	return true
}

func (s *CartServer) Reload() {}

func (s *CartServer) MainLoop() { time.Sleep(time.Second) }

func (s *CartServer) Final() bool { return true }

func main() {
	defer func() {
		glog.Flush()
	}()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	flag.Set("v", "2")
	flag.Parse()
	CartServerGetInstance().Main()
	glog.Infoln("[CartServer] server closed.")
}
