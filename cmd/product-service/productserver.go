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

type ProductServer struct {
	service.Service
}

var (
	server *ProductServer
	once   sync.Once
)

func ProductServerGetInstance() *ProductServer {
	once.Do(func() {
		server = &ProductServer{}
		server.Derived = server
	})
	return server
}

func (s *ProductServer) Init() bool {
	if !configs.ParseConfig() {
		glog.Errorln("[ProductServer] parse config error.")
		return false
	}
	if !rpcServerStart() {
		glog.Errorln("[ProductServer] rpc server start error.")
		return false
	}
	if !mysqlDatabaseInit() {
		glog.Errorln("[ProductServer] mysql database init error.")
		return false
	}
	// Consul register
	if !service.ServiceRegister(
		"1",
		"product-service",
		configs.GetConf().ProductCfg.Host,
		configs.GetConf().ProductCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[ProductServer] consul register error.")
		return false
	}
	return true
}

func (s *ProductServer) Reload() {
}

func (s *ProductServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *ProductServer) Final() bool {
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
	ProductServerGetInstance().Main()
	glog.Infoln("[ProductServer] server closed.")
}
