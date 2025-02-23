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

type PayServer struct {
	service.Service
}

var (
	server *PayServer
	once   sync.Once
)

func PayServerGetInstance() *PayServer {
	once.Do(func() {
		server = &PayServer{}
		server.Derived = server
	})
	return server
}

func (s *PayServer) Init() bool {
	if !configs.ParseConfig() {
		glog.Errorln("[PayServer] parse config error.")
		return false
	}
	if !rpcServerStart() {
		glog.Errorln("[PayServer] rpc server start error.")
		return false
	}
	if !mysqlDatabaseInit() {
		glog.Errorln("[PayServer] mysql database init error.")
		return false
	}
	// Consul register
	if !service.ServiceRegister(
		"1",
		"pay-service",
		configs.GetConf().PayCfg.Host,
		configs.GetConf().PayCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[PayServer] consul register error.")
		return false
	}
	return true
}

func (s *PayServer) Reload() {
}

func (s *PayServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *PayServer) Final() bool {
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
	PayServerGetInstance().Main()
	glog.Infoln("[PayServer] server closed.")
}
