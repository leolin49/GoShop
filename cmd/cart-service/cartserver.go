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

type CartServer struct {
	service.Service
}

var (
	server *CartServer
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
	if !configs.ParseConfig() {
		glog.Errorln("[CartServer] parse config error.")
		return false
	}
	if !rpcServerStart() {
		glog.Errorln("[CartServer] rpc server start error.")
		return false
	}
	if !mysqlDatabaseInit() {
		glog.Errorln("[CartServer] mysql database init error.")
		return false
	}
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

func (s *CartServer) Reload() {
}

func (s *CartServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *CartServer) Final() bool {
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
	CartServerGetInstance().Main()
	glog.Infoln("[CartServer] server closed.")
}
