package main

import (
	"flag"
	"goshop/configs"
	"goshop/pkg/redis"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

type CheckoutServer struct {
	service.Service
	rdb	*redis.Rdb
}

var (
	server *CheckoutServer
	once   sync.Once
)

func CheckoutServerGetInstance() *CheckoutServer {
	once.Do(func() {
		server = &CheckoutServer{}
		server.Derived = server
	})
	return server
}

func (s *CheckoutServer) Init() bool {
	if !configs.ParseConfig() {
		glog.Errorln("[CheckoutServer] parse config error.")
		return false
	}
	if !rpcServerStart() {
		glog.Errorln("[CheckoutServer] rpc server start error.")
		return false
	}

	// mysql
	if !mysqlDatabaseInit() {
		glog.Errorln("[CheckoutServer] mysql database init error.")
		return false
	}

	// redis
	s.rdb = redis.NewRedisClient(&configs.GetConf().RedisCfg)

	// rpc clients
	rpcClientsStart()

	// RabbitMQ consume
	go rabbitConsumer("checkout-queue", configs.GetConf().GetRabbitMQUrl())

	// Consul register
	if !service.ServiceRegister(
		"1",
		"checkout-service",
		configs.GetConf().CheckoutCfg.Host,
		configs.GetConf().CheckoutCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[CheckoutServer] consul register error.")
		return false
	}
	return true
}

func (s *CheckoutServer) Reload() {
}

func (s *CheckoutServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *CheckoutServer) Final() bool {
	return true
}

func main() {
	defer func() {
		rpcClientClose()
		glog.Flush()
	}()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	flag.Set("v", "2")
	flag.Parse()
	CheckoutServerGetInstance().Main()
	glog.Infoln("[CheckoutServer] server closed.")
}
