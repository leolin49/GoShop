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

type LoginServer struct {
	service.Service
}

var (
	server *LoginServer
	once   sync.Once
)

func LoginServerGetInstance() *LoginServer {
	once.Do(func() {
		server = &LoginServer{}
		server.Derived = server
	})
	return server
}

func (s *LoginServer) Init() bool {
	if !configs.ParseConfig() {
		glog.Errorln("[LoginServer] parse config error.")
		return false
	}
	if !rpcServerStart() {
		glog.Errorln("[LoginServer] rpc server start error.")
		return false
	}
	if !mysqlDatabaseInit() {
		glog.Errorln("[LoginServer] mysql database init error.")
		return false
	}

	rpcClientsStart()
	
	// Consul register
	if !service.ServiceRegister(
		"1",
		"login-service",
		configs.GetConf().LoginCfg.Host,
		configs.GetConf().LoginCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[LoginServer] consul register error.")
		return false
	}
	return true
}

func (s *LoginServer) Reload() {
}

func (s *LoginServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *LoginServer) Final() bool {
	return true
}

func main() {
	defer func() {
		AuthClientClose()
		glog.Flush()
	}()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	flag.Set("v", "2")
	flag.Parse()

	LoginServerGetInstance().Main()
	glog.Infoln("[LoginServer] server closed.")
}
