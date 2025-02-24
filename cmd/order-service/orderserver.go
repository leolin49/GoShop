package main

import (
	"flag"
	"goshop/configs"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

type OrderServer struct {
	service.Service
}

var (
	server *OrderServer
	once   sync.Once
)

func OrderServerGetInstance() *OrderServer {
	once.Do(func() {
		server = &OrderServer{}
		server.Derived = server
	})
	return server
}

func (s *OrderServer) Init() bool {
	if !configs.ParseConfig() {
		glog.Errorln("[OrderServer] parse config error.")
		return false
	}
	if !rpcServerStart() {
		glog.Errorln("[OrderServer] rpc server start error.")
		return false
	}
	if !mysqlDatabaseInit() {
		glog.Errorln("[OrderServer] mysql database init error.")
		return false
	}
	// Consul register
	if !service.ServiceRegister(
		"1",
		"order-service",
		configs.GetConf().OrderCfg.Host,
		configs.GetConf().OrderCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[OrderServer] consul register error.")
		return false
	}
	return true
}

func (s *OrderServer) Reload() {
}

func (s *OrderServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *OrderServer) Final() bool {
	return true
}

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
	OrderServerGetInstance().Main()
	glog.Infoln("[OrderServer] server closed.")
}
